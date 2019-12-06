## Advent 2019 log

This is a log of my experiences working on this years Advent of Code.  I never finish these, but it's always fun to get started and try.  I'm going to use Go this year rather than python.

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

### Day 3

If I was doing this in python I would just compute all the points on the lines and their manhattan distances as a single set per wire, then take the min value of the intersections of the sets.  I decided to do line segments and intersections in Go, which turned out to be more work.

For part one, I computed all the line segments and then compute all the intersections of the line segments in a nested loop.  I then loop over the intersections and compute their minimum manhattan distance, keeping track of the minimum distance.  Using segments is tricky because you need to make sure the order in which you handle points is consistent.  

For part two, I traverse the segments, keeping track of the distance travelled and then record that distance for each point/wire pair.  I use some `map`s to make it `O(1)` for whether there is a point in the intersection.

In hindsight, I solved this problem incorrectly.  I should have taken the "pixel" approach, given that I knew the board is discrete and the numbers of points aren't too big.  Instead I created a solution that is resolution independent which was complicated and time-consuming.