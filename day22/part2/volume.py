import collections

# The code here to compute the volume of overlapping cuboids is from:
# https://stackoverflow.com/questions/69137352/computing-the-volume-of-the-union-of-axis-aligned-cubes

Interval = collections.namedtuple("Interval", ["lower", "upper"])
Cube = collections.namedtuple("Cube", ["x", "y", "z"])

def length(interval):
    return interval.upper - interval.lower


def length_of_union(intervals):
    events = []
    for interval in intervals:
        events.append((interval.lower, 1))
        events.append((interval.upper, -1))
    events.sort()
    previous = None
    overlap = 0
    total = 0
    for x, delta in events:
        if overlap > 0:
            total += x - previous
        previous = x
        overlap += delta
    return total


def all_boundaries(intervals):
    boundaries = set()
    for interval in intervals:
        boundaries.add(interval.lower)
        boundaries.add(interval.upper)
    return sorted(boundaries)


def subdivide_at(interval, boundaries):
    lower = interval.lower
    for x in sorted(boundaries):  # Resort is O(n) due to Timsort.
        if x < lower:
            pass
        elif x < interval.upper:
            yield Interval(lower, x)
            lower = x
        else:
            yield Interval(lower, interval.upper)
            break


def volume_of_union(cubes):
    x_boundaries = all_boundaries(cube.x for cube in cubes)
    y_boundaries = all_boundaries(cube.y for cube in cubes)
    sub_problems = collections.defaultdict(list)
    for cube in cubes:
        for x in subdivide_at(cube.x, x_boundaries):
            for y in subdivide_at(cube.y, y_boundaries):
                sub_problems[(x, y)].append(cube.z)
    return sum(
        length(x) * length(y) * length_of_union(z_intervals)
        for ((x, y), z_intervals) in sub_problems.items()
    )


# Note: My code is below

with open('./new-input.txt') as f:
    lines = f.read().splitlines()

cubes = [] 
for line in lines:
    parts = line.split()
    ranges = parts[1].split(',')
    xRange = ranges[0].split('=')[1].split('..')
    yRange = ranges[1].split('=')[1].split('..')
    zRange = ranges[2].split('=')[1].split('..')
    cube = Cube(Interval(int(xRange[0]),int(xRange[1])+1), Interval(int(yRange[0]),int(yRange[1])+1), Interval(int(zRange[0]),int(zRange[1])+1))
    cubes.append(cube) 
    
print(volume_of_union(cubes))