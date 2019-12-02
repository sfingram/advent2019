## Advent 2019 log

### Day 1

#### Part 1

Trivial: iterate over lines in a file and call a function on each while maintaining a running sum.

#### Part 2

The main trick here is that the exercise wants you to perform the recursive fuel calculation on each mass individually.  The sum of the recursive fuel calculations is not equivalent to the recursive fuel calculation of the original fuel sum.

### Day 2

#### Part 1

I initially errored by thinking that the input data was line-by-line instead of comma separated on a single line.  I modified my scanner with a different comma-based split function (one I found in a test file somewhere).  The next issue I had was an off-by-one error, forgetting that the `for` loop will automatically increment the instruction pointer by 1.

#### Part 2

Went very smoothly.  Refactored the previous solution to be paramterized.  Just doing a grid search over the parameter space was enough to efficiently find the answer.

Did a little post-cleanup to reduce the number of allocations instead of allocating inside of the loop.