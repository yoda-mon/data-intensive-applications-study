/*
 * https://doc.akka.io/docs/akka/current/typed/actors.html#first-example
 */


import akka.actor.typed.scaladsl.Behaviors
import akka.actor.typed.scaladsl.LoggerOps
import akka.actor.typed.{ ActorRef, ActorSystem, Behavior }

// あいさつしろ命令を受けたらあいさつし、あいさつ確認を送るアクター
object HelloWorld {
  // 2種類のメッセージを定義
  final case class Greet(whom: String, replyTo: ActorRef[Greeted])  // 誰かにあいさつしろという命令メッセージ whom:あいさつ相手 replyTo: あいさつ確認を送る相手
  final case class Greeted(whom: String, from: ActorRef[Greet])     // あいさつ確認をするメッセージ whom: あいさつした相手 from: あいさつした人

  def apply(): Behavior[Greet] = Behaviors.receive { (context, message) =>  // Greetメッセージをreceiveしたときの挙動を定義
    context.log.info("Hello {}!", message.whom)  // あいさつ
    println(s"Hello ${message.whom}!")
    message.replyTo ! Greeted(message.whom, context.self)  // replyToに向かってGreetedメッセージを送信
    Behaviors.same
  }
}

// あいさつ確認を受けたら再度あいさつ命令を送り返すアクター
object HelloWorldBot {

  def apply(max: Int): Behavior[HelloWorld.Greeted] = {  // Greetedを送ったときの挙動を定義
    bot(0, max)
  }

  private def bot(greetingCounter: Int, max: Int): Behavior[HelloWorld.Greeted] =
    Behaviors.receive { (context, message) =>
      val n = greetingCounter + 1                                       // カウンタアップ
      context.log.info2("Greeting {} for {}", n, message.whom)
      println(s"Greeting $n for ${message.whom}")
      if (n == max) {
        Behaviors.stopped
      } else {
        message.from ! HelloWorld.Greet(message.whom, context.self)    // Greetedしたアクターに向かってGreet
        bot(n, max)                                                    // 再帰
      }
    }
}

// 全体を統括するアクター
// SayHelloメッセージを受けたら上2つのアクターを生成しあいさつさせる
object HelloWorldMain {
  // SayHelloメッセージを定義
  final case class SayHello(name: String)

  def apply(): Behavior[SayHello] =
    Behaviors.setup { context =>
      val greeter = context.spawn(HelloWorld(), "greeter")  // HelloWorldアクター起動

      Behaviors.receiveMessage { message =>
        val replyTo = context.spawn(HelloWorldBot(max = 3), message.name)  // HelloWorldBotアクター起動
        greeter ! HelloWorld.Greet(message.name, replyTo)  // 起動したHelloWorldアクターに、HelloWorldBotに向かってあいさつしてからあいさつ確認しろとメッセージ
        Behaviors.same
      }
    }

  def main(args: Array[String]): Unit = {
    val system: ActorSystem[HelloWorldMain.SayHello] =  // HelloWorldMainアクターを起動
      ActorSystem(HelloWorldMain(), "hello")

    system ! HelloWorldMain.SayHello("World")  // WorldとあいさつしろとSayHelloメッセージ
    system ! HelloWorldMain.SayHello("Akka")    // AkkaとあいさつしろとSayHelloメッセージ
  }
}

