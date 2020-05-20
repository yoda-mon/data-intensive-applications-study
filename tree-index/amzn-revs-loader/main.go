package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	//"encoding/csv"
	"database/sql"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"github.com/yoda-mon/data-intensive-application/tree-index/amzn-revs-loader/csv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initZapLog() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	//config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	logger, _ := config.Build()
	return logger
}

var logger *zap.SugaredLogger

func init() {
	loggerMgr := initZapLog()
	zap.ReplaceGlobals(loggerMgr)
	defer loggerMgr.Sync() // flushes buffer, if any
	logger = loggerMgr.Sugar()
}

type Cell interface {
	GetStrValue() string
	IsValid()
}

type Record struct {
	cells [15]Cell // Arrayのほうがポインタ処理入らなくて良さげ
}

func (r *Record) setCells(e []string, schema []string) {
	var re = regexp.MustCompile(`\((.*?)\)`)
	for i, s := range schema {
		logger.Debug("Column: ", i)
		switch t := s; {
		case strings.HasPrefix(t, "varchar"):
			l, _ := strconv.Atoi(re.FindStringSubmatch(t)[1])
			r.cells[i] = &VarCharCell{value: e[i], length: l}
			logger.Debug(r.cells[i].GetStrValue())
		case t == "integer":
			r.cells[i] = &IntCell{value: e[i]}
			logger.Debug(r.cells[i].GetStrValue())
		case t == "date":
			r.cells[i] = &DateCell{value: e[i]}
			logger.Debug(r.cells[i].GetStrValue())
		case t == "text":
			r.cells[i] = &TextCell{value: e[i]}
			logger.Debug(r.cells[i].GetStrValue())
		default:
			r.cells[i] = &VarCharCell{value: "", length: 0}
			logger.Warn("Not defined type:")
		}
	}
}

func (r *Record) validateCells() {
	logger.Debug(len(r.cells))
	for _, v := range r.cells {
		v.IsValid()
	}
}

func main() {

	//file, err := os.Open("../data/sample_file.tsv")
	file, err := os.Open("../data/amazon_reviews_multilingual_JP_v1_00.tsv")
	//file, err := os.Open("../data/error_samples.tsv")
	if err != nil {
		logger.Error(err)
	}
	defer file.Close()

	connStr := "user=student dbname=student sslmode=disable host=localhost port=5432 password=1234"
	/*
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}
			row, err := db.Query("SELECT * FROM amzn_revs.multailingual_jp")
			defer row.Close()
			for row.Next() {
				var marketplace string
				var customerID int
				var reviewID string
				var productID string
				var productParent int
				var productTitle string
				var productCategory string
				var starRating int
				var helpfulVotes int
				var totalVotes int
				var vine string
				var verifiedPurchase string
				var reviewHeadline string
				var reviewBody string
				var reviewDate string
				if err := row.Scan(
					&marketplace,
					&customerID,
					&reviewID,
					&productID,
					&productParent,
					&productTitle,
					&productCategory,
					&starRating,
					&helpfulVotes,
					&totalVotes,
					&vine,
					&verifiedPurchase,
					&reviewHeadline,
					&reviewBody,
					&reviewDate,
				); err != nil {
					log.Fatal(err)
				}
				logger.Debug(marketplace, customerID)
			}
	*/
	ch := make(chan Record, 1000)
	numRoutines := 4
	var wg sync.WaitGroup
	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func() {
			db, err := sql.Open("postgres", connStr) // TODO: スレッドセーフならループ外だし
			if err != nil {
				log.Fatal(err)
			}
			for {
				r, ok := <-ch
				if !ok {
					wg.Done()
					return
				}
				sql := fmt.Sprintf(`
				INSERT INTO amzn_revs.multailingual_jp
				(
					marketplace,
					customer_id,
					review_id,
					product_id,
					product_parent,
					product_title,
					product_category,
					star_rating,
					helpful_votes,
					total_votes,
					vine,
					verified_purchase,
					review_headline,
					review_body,
					review_date 
				)
				VALUES ('%s', %s, '%s', '%s', %s, %s, '%s', %s, %s, %s, '%s', '%s', %s, %s, '%s')
				`,
					r.cells[0].GetStrValue(),
					r.cells[1].GetStrValue(),
					r.cells[2].GetStrValue(),
					r.cells[3].GetStrValue(),
					r.cells[4].GetStrValue(),
					pq.QuoteLiteral(r.cells[5].GetStrValue()),
					r.cells[6].GetStrValue(),
					r.cells[7].GetStrValue(),
					r.cells[8].GetStrValue(),
					r.cells[9].GetStrValue(),
					r.cells[10].GetStrValue(),
					r.cells[11].GetStrValue(),
					pq.QuoteLiteral(r.cells[12].GetStrValue()),
					pq.QuoteLiteral(r.cells[13].GetStrValue()),
					r.cells[14].GetStrValue(),
				)
				logger.Debug(sql)
				_, err := db.Exec(sql)
				if err != nil {
					logger.Error(sql)
					logger.Errorf("%v", err)
				}
			}
		}()
	}

	//r := csv.NewReader(bufio.NewReader(file))
	r := csv.NewReader(bufio.NewReader(file))
	r.Comma = '\t'             // Field delimiter
	r.Comment = '#'            // Comment charactor
	r.FieldsPerRecord = 15     // Number of expected fields per record.
	r.LazyQuotes = false       // Allow quote in a unquoted field and non-quote in a quoted field
	r.TrimLeadingSpace = false // Ignore leading white space in a feld
	r.ReuseRecord = false      // Reuse a slice of previous read call to improve performance

	_, _ = r.Read() // Skip header row

	schema := [15]string{
		"varchar(2)",   // marketplace
		"integer",      // customer_id
		"varchar(14)",  // reveiw_id
		"varchar(10)",  // product_id
		"integer",      // product_parent
		"varchar(400)", // product_title
		"varchar(40)",  // product_category
		"integer",      // star_rating
		"integer",      // helpful_votes
		"integer",      // total_votes
		"varchar(1)",   // vine
		"varchar(1)",   // verified_purchase
		"varchar(300)", // review_headline
		"text",         // review_body
		"date",         // review_date
	}

	lineNum := 1
	errRows := 0
	for {
		lineNum++
		logger.Debug("Line number: ", lineNum)
		elms, err := r.Read()
		if err == io.EOF {
			logger.Warn("EOF")
			break
		} else if err1, ok := err.(*csv.ParseError); ok && err1.Err == csv.ErrFieldCount {
			//logger.Error("Error field count on line: ", err1.Line)
			logger.Error(err1)
			errRows++
			continue
		} else if err2, ok := err.(*csv.ParseError); ok && err2.Err == csv.ErrQuote {
			//logger.Warn("Extraneous or missing double-quote on line: ", err2.Line) // r.LazyQuotesをfalseにしてればこちらに落ちて継続してくれる
			logger.Warn(err2)
			errRows++
			continue
		} else if err3, ok := err.(*csv.ParseError); ok && err3.Err == csv.ErrBareQuote {
			//logger.Warn("Bare double-quote on line: ", err3.Line) // r.LazyQuotesをfalseにしてればこちらに落ちて継続してくれる
			logger.Warn(err3)
			errRows++
			continue
		} else if err != nil {
			logger.Error(err)
		} else {
			record := Record{}
			record.setCells(elms, schema[:])

			record.validateCells()
			ch <- record
		}
		/*

			var marketplace Cell = VarCharCell{value: record[0], length: 2}
			marketplace.IsValid()
			logger.Info(marketplace.GetStrValue())

			var customerID Cell = IntCell{value: record[1]}
			customerID.IsValid()
			logger.Info(customerID.GetStrValue())

			var reviewDate Cell = DateCell{value: record[14]}
			reviewDate.IsValid()
			logger.Info(reviewDate.GetStrValue())
		*/
		/*
			r := <-ch
			logger.Debug(r.cells[0].GetStrValue())
		*/
	}
	close(ch)
	wg.Wait()
	logger.Warn("Error rows: ", errRows)
}
