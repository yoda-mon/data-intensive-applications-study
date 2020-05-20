package main;

import java.nio.file.Paths;
import java.nio.file.Files;
import java.io.File;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.dataformat.avro.AvroMapper;
import com.fasterxml.jackson.dataformat.protobuf.ProtobufFactory;
import com.fasterxml.jackson.dataformat.protobuf.schema.ProtobufSchemaLoader;
import org.msgpack.jackson.dataformat.MessagePackFactory;

import thrift.Person;
import org.apache.thrift.TSerializer;
import org.apache.thrift.protocol.TBinaryProtocol;


public class Main {
    public static void main(String[] args) {
        var person = new main.Person();
        var jsonObjectMapper = new ObjectMapper();
        /*
        JSON to Java object
         */
        try {
            /*
            var json = "{\"userName\":\"Martin\",\"favoriteNumber\":1337,\"interests\":[\"daydreaming\",\"hacking\"]}";
            sample = jsonObjectMapper.readValue(json, Sample.class);
             */
            person = jsonObjectMapper.readValue(Paths.get("sample.json").toFile(), main.Person.class);
            var data = jsonObjectMapper.writeValueAsBytes(person);
            Files.write(new File("dist/sample.min.json").toPath(), data);  // minified json
        } catch (Exception ex) {
            ex.printStackTrace();
        }
        System.out.println(person);


        /*
        MessagePack
        Java object to mpac
         */
        var messagePackObjectMapper = new ObjectMapper(new MessagePackFactory());
        try {
            var data = messagePackObjectMapper.writeValueAsBytes(person);
            Files.write(new File("dist/sample.mpac").toPath(), data);
        } catch (Exception ex) {
            ex.printStackTrace();
        }
        /*
        Apache thrift
        Java object to thrift
         */
        try {
            var personThrift = new thrift.Person();
            personThrift = jsonObjectMapper.readValue(Paths.get("sample.json").toFile(), thrift.Person.class);
            var serializer = new TSerializer(new TBinaryProtocol.Factory());
            var data = serializer.serialize(personThrift);
            Files.write(new File("dist/sample.tb").toPath(), data);  // 拡張子名は？
        } catch (Exception ex) {
            ex.printStackTrace();
        }

        /*
        Protocol Buffer
        Java object to protobuf
         */
        var protobufObjectMapper = new ObjectMapper(new ProtobufFactory());
        try {
            var schema = ProtobufSchemaLoader.std.load(Paths.get("config/sample.proto").toFile());
            var data = protobufObjectMapper.writer(schema).writeValueAsBytes(person);
            Files.write(new File("dist/sample.pb").toPath(), data);  // 拡張子名は？
        } catch (Exception ex) {
            ex.printStackTrace();
        }

        /*
        Apache Avro
        Java object to avro
         */
        var avroObjectMapper = new AvroMapper();
        try {
            var schema = avroObjectMapper.schemaFrom(Paths.get("config/sample.avsc").toFile());
            var data = avroObjectMapper.writer(schema).writeValueAsBytes(person);
            Files.write(new File("dist/sample.ab").toPath(), data);  // 拡張子名は？
        } catch (Exception ex) {
            ex.printStackTrace();
        }

    }
}
