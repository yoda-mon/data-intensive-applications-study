namespace java thrift
namespace go thrift

/*
 C like comments are supported
*/
// This is also a valid comment

typedef i64 int // We can use typedef to get pretty names for the types we are using
struct Person {
    1: required string userName,
    2: optional int favoriteNumber,
    3: optional list<string> interests
}
