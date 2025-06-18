use std::{cmp, str};

fn main() {
    let input = include_str!("./input.txt");
    println!("{}", find_seed_with_neares_location(input));
}

#[derive(Clone, Debug, PartialEq)]
struct Range {
    start: i64,
    len: i64,
}

#[derive(Clone, Debug, PartialEq)]
struct PreRemapRange {
    start: i64,
    len: i64,
    remap_offset: i64,
}

#[derive(Debug)]
struct MapEntry {
    dst_start: i64,
    src_start: i64,
    len: i64,
}

fn find_seed_with_neares_location(input: &str) -> i64 {
    let mut iter = input.lines();

    let seed_line = iter.next().unwrap();
    // println!("seed_line = '{}'", seed_line);
    let seed_nums = seed_line["seeds: ".len()..]
        .split(' ')
        .map(|s| s.parse::<i64>().unwrap())
        .collect::<Vec<_>>();
    // println!("seed_nums = {:?}", seed_nums);

    let seed_ranges = seed_nums
        .chunks(2)
        .map(|x| Range {
            start: x[0],
            len: x[1],
        })
        .collect::<Vec<_>>();
    println!("seed_ranges = {:?}", seed_ranges);

    let mut seed_to_soil: Vec<MapEntry> = vec![];
    let mut soil_to_fertilizer: Vec<MapEntry> = vec![];
    let mut fertilizer_to_water: Vec<MapEntry> = vec![];
    let mut water_to_light: Vec<MapEntry> = vec![];
    let mut light_to_temperature: Vec<MapEntry> = vec![];
    let mut temperature_to_humidity: Vec<MapEntry> = vec![];
    let mut humidity_to_location: Vec<MapEntry> = vec![];

    while let Some(l) = iter.next() {
        match l {
            "seed-to-soil map:" => seed_to_soil = read_map(&mut iter),
            "soil-to-fertilizer map:" => soil_to_fertilizer = read_map(&mut iter),
            "fertilizer-to-water map:" => fertilizer_to_water = read_map(&mut iter),
            "water-to-light map:" => water_to_light = read_map(&mut iter),
            "light-to-temperature map:" => light_to_temperature = read_map(&mut iter),
            "temperature-to-humidity map:" => temperature_to_humidity = read_map(&mut iter),
            "humidity-to-location map:" => humidity_to_location = read_map(&mut iter),
            _ => println!("Line with unexpected value '{}'", l),
        }
    }
    // println!("seed-to-soil map = {:?}", seed_to_soil);
    // println!("soil_to_fertilizer map = {:?}", soil_to_fertilizer);
    // println!("fertilizer_to_water map = {:?}", fertilizer_to_water);
    // println!("water_to_light map = {:?}", water_to_light);
    // println!("light_to_temperature map = {:?}", light_to_temperature);
    // println!(
    //     "temperature_to_humidity map = {:?}",
    //     temperature_to_humidity
    // );
    // println!("humidity_to_location map = {:?}", humidity_to_location);

    let soil_ranges = intersect_remap_ranges(&seed_ranges, &seed_to_soil);
    let fertilizer_ranges = intersect_remap_ranges(&soil_ranges, &soil_to_fertilizer);
    let water_ranges = intersect_remap_ranges(&fertilizer_ranges, &fertilizer_to_water);
    let light_ranges = intersect_remap_ranges(&water_ranges, &water_to_light);
    let temperature_ranges = intersect_remap_ranges(&light_ranges, &light_to_temperature);
    let humidity_ranges = intersect_remap_ranges(&temperature_ranges, &temperature_to_humidity);
    let location_ranges = intersect_remap_ranges(&humidity_ranges, &humidity_to_location);

    // println!("soil_ranges = {:?}", soil_ranges);
    // println!("fertilizer_ranges = {:?}", fertilizer_ranges);
    // println!("water_ranges = {:?}", water_ranges);
    // println!("light_ranges = {:?}", light_ranges);
    // println!("temperature_ranges = {:?}", temperature_ranges);
    // println!("humidity_ranges = {:?}", humidity_ranges);
    // println!("location_ranges = {:?}", location_ranges);

    location_ranges.iter().map(|r| r.start).min().unwrap()
}

fn read_map(iter: &mut str::Lines) -> Vec<MapEntry> {
    let mut map: Vec<MapEntry> = vec![];
    for ll in iter.by_ref() {
        if ll.is_empty() {
            break;
        }
        let nums: Vec<i64> = ll.split(' ').map(|s| s.parse::<i64>().unwrap()).collect();
        map.push(MapEntry {
            dst_start: nums[0],
            src_start: nums[1],
            len: nums[2],
        });
    }
    map
}

fn intersect_remap_ranges(ranges: &[Range], map: &[MapEntry]) -> Vec<Range> {
    let mut pre_remap_ranges = ranges
        .iter()
        .map(|r| PreRemapRange {
            start: r.start,
            len: r.len,
            remap_offset: 0,
        })
        .collect::<Vec<_>>();

    for m in map.iter() {
        pre_remap_ranges = intersect_pre_remap_ranges_with_entry(&pre_remap_ranges, m);
    }

    pre_remap_ranges
        .iter()
        .map(|pr| Range {
            start: pr.start + pr.remap_offset,
            len: pr.len,
        })
        .collect::<Vec<_>>()
}

fn intersect_pre_remap_ranges_with_entry(
    ranges: &[PreRemapRange],
    map_entry: &MapEntry,
) -> Vec<PreRemapRange> {
    ranges
        .iter()
        .flat_map(|r| intersect_pre_remap_range_with_entry(r, map_entry))
        .collect::<Vec<_>>()
}

fn intersect_remap_range_with_entry(range: &Range, map_entry: &MapEntry) -> Vec<Range> {
    intersect_range_with_entry(range, map_entry)
        .iter()
        .map(|pr| Range {
            start: pr.start + pr.remap_offset,
            len: pr.len,
        })
        .collect::<Vec<_>>()
}

fn intersect_range_with_entry(range: &Range, map_entry: &MapEntry) -> Vec<PreRemapRange> {
    intersect_pre_remap_range_with_entry(
        &PreRemapRange {
            start: range.start,
            len: range.len,
            remap_offset: 0,
        },
        map_entry,
    )
}

fn intersect_pre_remap_range_with_entry(
    range: &PreRemapRange,
    map_entry: &MapEntry,
) -> Vec<PreRemapRange> {
    /*
    // No intersection?
    if range.start >= map_entry.src_start + map_entry.len {
        return vec![range];
    }
    // No intersection?
    if map_entry.src_start >= range.start + range.len {
        return vec![range];
    }
    */

    let remap_offset = map_entry.dst_start - map_entry.src_start;

    // Map range fully within range?
    if map_entry.src_start >= range.start
        && map_entry.src_start + map_entry.len <= range.start + range.len
    {
        let mut res = vec![];

        // Left range
        let left_len = map_entry.src_start - range.start;
        if left_len > 0 {
            res.push(PreRemapRange {
                start: range.start,
                len: left_len,
                remap_offset: range.remap_offset,
            });
        }

        // Remap center range
        res.push(PreRemapRange {
            start: map_entry.src_start, // Remap
            len: map_entry.len,
            remap_offset,
        });

        // Right range
        let right_len = (range.start + range.len) - (map_entry.src_start + map_entry.len);
        if right_len > 0 {
            res.push(PreRemapRange {
                start: map_entry.src_start + map_entry.len,
                len: right_len,
                remap_offset: range.remap_offset,
            });
        }

        return res;
    }

    // Range fully within map range?
    if range.start >= map_entry.src_start
        && (range.start + range.len) <= (map_entry.src_start + map_entry.len)
    {
        // Remap center range
        let offset = range.start - map_entry.src_start;
        return vec![PreRemapRange {
            start: map_entry.src_start + offset, // Remap
            len: range.len,
            remap_offset,
        }];
    }

    // Map range ends within range
    if (map_entry.src_start + map_entry.len) > range.start
        && (map_entry.src_start + map_entry.len) < (range.start + range.len)
    {
        let mut res = vec![];

        // Remap the left part of the range
        let left_len = (map_entry.src_start + map_entry.len) - range.start;
        let offset = map_entry.len - left_len;
        res.push(PreRemapRange {
            start: map_entry.src_start + offset,
            len: left_len,
            remap_offset,
        });

        // Right part of range stays unchanged
        let right_len = range.len - left_len;
        res.push(PreRemapRange {
            start: range.start + left_len,
            len: right_len,
            remap_offset: range.remap_offset,
        });

        return res;
    }

    // Map range starts within range
    if map_entry.src_start > range.start && map_entry.src_start < (range.start + range.len) {
        let mut res = vec![];

        let right_len = (range.start + range.len) - map_entry.src_start;

        // Left part stays unchanged
        let left_len = range.len - right_len;
        res.push(PreRemapRange {
            start: range.start,
            len: left_len,
            remap_offset: range.remap_offset,
        });

        // Remap right part of intersection
        res.push(PreRemapRange {
            start: map_entry.src_start,
            len: right_len,
            remap_offset,
        });

        return res;
    }

    vec![PreRemapRange {
        start: range.start,
        len: range.len,
        remap_offset: range.remap_offset,
    }]
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_find_seed_with_neares_location() {
        let input = concat!(
            r"seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4"
        );
        assert_eq!(find_seed_with_neares_location(input), 46);
    }

    #[test]
    fn test_intersect_remap() {
        // No intersection
        assert_eq!(
            intersect_remap_range_with_entry(
                &Range { start: 2, len: 2 },
                &MapEntry {
                    dst_start: 10,
                    src_start: 0,
                    len: 2
                }
            ),
            vec![Range { start: 2, len: 2 }]
        );
        // No intersection
        assert_eq!(
            intersect_remap_range_with_entry(
                &Range { start: 2, len: 2 },
                &MapEntry {
                    dst_start: 15,
                    src_start: 5,
                    len: 2
                }
            ),
            vec![Range { start: 2, len: 2 }]
        );
        // Map range fully within range
        assert_eq!(
            intersect_remap_range_with_entry(
                &Range { start: 1, len: 3 },
                &MapEntry {
                    dst_start: 10,
                    src_start: 2,
                    len: 1
                }
            ),
            vec![
                Range { start: 1, len: 1 },  // No remap
                Range { start: 10, len: 1 }, // Remapped
                Range { start: 3, len: 1 }   // No remap
            ]
        );
        // Range fully within map range
        assert_eq!(
            intersect_remap_range_with_entry(
                &Range { start: 2, len: 1 },
                &MapEntry {
                    dst_start: 10,
                    src_start: 1,
                    len: 3
                }
            ),
            vec![
                Range { start: 11, len: 1 }, // Remapped
            ]
        );
        // Map range ends within range
        assert_eq!(
            intersect_remap_range_with_entry(
                &Range { start: 1, len: 3 },
                &MapEntry {
                    dst_start: 10,
                    src_start: 0,
                    len: 3
                }
            ),
            vec![
                Range { start: 11, len: 2 }, // Remapped
                Range { start: 3, len: 1 },  // No remap
            ]
        );
        // Map range starts within range
        assert_eq!(
            intersect_remap_range_with_entry(
                &Range { start: 1, len: 3 },
                &MapEntry {
                    dst_start: 10,
                    src_start: 2,
                    len: 3
                }
            ),
            vec![
                Range { start: 1, len: 1 },  // No remap
                Range { start: 10, len: 2 }, // Remap
            ]
        );
    }
}
