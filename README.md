# remove-duplicate-files

A script I wrote for personal use, to delete duplicate files in a specific folders preferring some folders.

I wrote it initially in TypeScript but for speed purposes rewrote it in golang, it's also the first time ever I use golang.

Takeaways:
- More threads doesn't mean faster execution when the bottleneck is in the Hard Drive.
- I got to use some threading knowledge that I didn't get to use in other languages before (Mutexes and Wait Groups).
- Golang is conscise not as I was expecting and its standard library is good.
