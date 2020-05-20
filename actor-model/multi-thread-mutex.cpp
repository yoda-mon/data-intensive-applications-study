#include<iostream>
#include <thread>
#include <mutex>

void count_up(int* count, std::mutex* mtx)
{
    for(int i=0; i<1'000'000; i++){
        std::lock_guard<std::mutex> lock(*mtx);  // mutexをロック
        (*count)++;
    }
}

int main()
{
    int count = 0;
    std::mutex mtx; //これが他のスレッドににロックされていると先に進まずブロックされる

    std::thread t1(count_up, &count, &mtx);
    std::thread t2(count_up, &count, &mtx);

    t1.join();
    t2.join();

    std::cout << count << std::endl;

    return 0;
}