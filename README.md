# Threaded Random Tests

Originally, this project started as wanting to take a bogo sort and thread it to make as many shuffles as fast as possible. After developing a first round, I found that with each passing thread, I lost a slough of performance. I created this thread to investigate.

## Hypothesis one
bogo sort relies on the rand c lib function. There may be a chance that this is a shared resource and if a single thread was calling this function, it would be able to optimize, but since two threads have to share the resource, one would hold onto it while the other waits, effectively slowing down both threads.

As I continued to work on this, I found that removing rand() did speed up the threads, but my threads were still slower than expected. In two threads with no data dependencies and no locks, I'd expect the system to be able to perform N\*threads iterations, but the number of iterations didn't change between the number of threads. Two threads just took twice as long with twice as many iterations. So what's the deal? Why no speedup.

## Rabbit Hole
From here I found that if I used the struct to hold all my information and access and wrote to that, I wouldn't receive any speedup from additional threads -- despite each thread having it's own independant struct. Instead, if I create a local variable and transfer the info at the very end, I would get the N\*threads iterations I would expect. What the hell?!
