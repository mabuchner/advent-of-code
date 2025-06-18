fn main() {
    let input = include_str!("./input.txt");
    println!("{}", calc_product_of_win_possibilities(input));
}

fn calc_product_of_win_possibilities(input: &str) -> i32 {
    let mut lines = input.lines();

    let time_line = lines.next().unwrap();
    // println!("time_line = '{}'", time_line);
    let times = time_line
        .strip_prefix("Time:")
        .unwrap()
        .split(' ')
        .filter(|s| !s.is_empty())
        .map(|s| s.parse::<i32>().unwrap())
        .collect::<Vec<_>>();
    // println!("times = {:?}", times);

    let distance_line = lines.next().unwrap();
    // println!("distance_line = '{}'", distance_line);
    let distances = distance_line
        .strip_prefix("Distance:")
        .unwrap()
        .split(' ')
        .filter(|s| !s.is_empty())
        .map(|s| s.parse::<i32>().unwrap())
        .collect::<Vec<_>>();
    // println!("distances = {:?}", distances);

    if times.len() != distances.len() {
        panic!(
            "Mismatchin sizes, times.len() = {}, distances.len() = {}",
            times.len(),
            distances.len()
        );
    }

    times
        .iter()
        .zip(distances.iter())
        .map(|(race_time, win_distance)| {
            let mut win_count = 0;
            for press_time in 1..*race_time {
                let speed = press_time; // mm/ms
                let drive_time = race_time - press_time; // ms
                let distance = speed * drive_time;
                if distance > *win_distance {
                    win_count += 1;
                }
            }
            win_count
        })
        .reduce(|prod, win_count| prod * win_count)
        .unwrap()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_calc_product_of_win_possibilities() {
        let input = concat!("Time:      7  15   30\n", "Distance:  9  40  200\n",);
        assert_eq!(calc_product_of_win_possibilities(input), 288);
    }
}
