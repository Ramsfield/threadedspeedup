#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <pthread.h>

#define RANDS 3117410583
#define NUMTHREADS 4
#define release(t) do{if(t) {free(t); t = NULL;}} while(0)


typedef struct {
  unsigned long randnums;
  int rnum;
} Info;

void * rfunc(void * vinfo) {
  Info * info = (Info *)vinfo;
#ifndef SLOWDOWN
  int rnum = 0;
  for(unsigned long tries = 0; tries < RANDS; ++tries) {
    rnum += 8000 % 4000;
  }
  info->rnum = rnum;
#else
  int rnum = 0;
  for(unsigned long tries = 0; tries < RANDS; ++tries) {
    info->rnum += 8000 % 4000;
  }
#endif
  info->randnums += RANDS;
  return NULL;
}

double getTimeElapsed(const struct timespec * start, const struct timespec * end) {
  double s,e;
  s = start->tv_sec + (start->tv_nsec / 1e+9);
  e = end->tv_sec + (end->tv_nsec / 1e+9);
  return e-s;
}

int main(int argc, char * argv[]) {
  //Seed random
  srand(time(NULL));

  //Variables
  int numthreads = NUMTHREADS;
  pthread_t * threads = NULL;
  Info * tinfo = NULL;
  struct timespec start,end;
  int rc;
  unsigned long total_rands = 0;

  //Functionality
  if(argc > 1) {
    numthreads = atoi(argv[1]);
  }
  threads = (pthread_t *) malloc(sizeof(pthread_t) * numthreads);
  tinfo = (Info *) malloc(sizeof(Info) * numthreads);
  clock_gettime(CLOCK_REALTIME, &start);
  for(int i=0; i < numthreads; ++i) {
    tinfo[i].randnums = 0;
    tinfo[i].rnum = 0;
    rc = pthread_create(threads + i, NULL, rfunc, (void*)(tinfo + i));
    if(rc) {
      printf("Error; return code from pthread_create() is %d\n", rc);
      goto exit;
    }
  }

  for(int i=0; i < numthreads; ++i) {
    pthread_join(threads[i], NULL);
    total_rands += tinfo[i].randnums;
  }

  clock_gettime(CLOCK_REALTIME, &end);

  double tElapsed = getTimeElapsed(&start, &end);
  printf("%lu randoms in %f seconds over %d threads. %f randoms per second\n",
      total_rands,
      tElapsed,
      numthreads,
      total_rands / tElapsed);

exit:
  release(threads);
  release(tinfo);
  return 0;
}
