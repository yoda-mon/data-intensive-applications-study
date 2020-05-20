#include<iostream>
#include <thread>


void count_up(int* count)
{
    for(int i=0; i<1'000'000; i++){
        (*count)++;
    }
}

int main()
{
    int count = 0;

    std::thread t1(count_up, &count);
    std::thread t2(count_up, &count);

    t1.join();
    t2.join();

    std::cout << count << std::endl;

    return 0;
}