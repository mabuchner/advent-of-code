use std::{cmp, str};

fn main() {
    let input = include_str!("./input.txt");
    println!("{}", find_seed_with_neares_location(input));
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
    let seeds: Vec<i64> = seed_line["seeds: ".len()..]
        .split(' ')
        .map(|s| s.parse::<i64>().unwrap())
        .collect();
    // println!("seeds = {:?}", seeds);

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

    let mut min_location = i64::MAX;
    for seed in seeds {
        let soil = lookup(&seed_to_soil, seed);
        let fertilizer = lookup(&soil_to_fertilizer, soil);
        let water = lookup(&fertilizer_to_water, fertilizer);
        let light = lookup(&water_to_light, water);
        let temperature = lookup(&light_to_temperature, light);
        let humidity = lookup(&temperature_to_humidity, temperature);
        let location = lookup(&humidity_to_location, humidity);
        min_location = cmp::min(min_location, location);
    }
    min_location
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

fn lookup(map: &Vec<MapEntry>, src: i64) -> i64 {
    for entry in map {
        if src >= entry.src_start && src < entry.src_start + entry.len {
            let offset = src - entry.src_start;
            return entry.dst_start + offset;
        }
    }
    src
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
        assert_eq!(find_seed_with_neares_location(input), 35);
    }
}
