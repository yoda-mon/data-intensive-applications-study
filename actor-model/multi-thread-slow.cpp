#include<iostream>
#include <thread>
#include <mutex>

void count_up(int* count, std::mutex* mtx)
{
    for(int i=0; i<1'000'000; i++){
        std::lock_guard<std::mutex> lock(*mtx);
        (*count)++;
    }
}

void count_up_slowly(int* count, std::mutex* mtx)
{
    for(int i=0; i<1'000'000; i++){
        std::lock_guard<std::mutex> lock(*mtx);
        std::this_thread::sleep_for(std::chrono::microseconds(1));
        (*count)++;
    }
}

int main()
{
    int count = 0;
    std::mutex mtx; //これが他のスレッドににロックされていると先に進まずブロックされる

    std::thread t1(count_up, &count, &mtx);
    std::thread t2(count_up_slowly, &count, &mtx);

    t1.join();
    t2.join();

    std::cout << count << std::endl;

    return 0;
}