fn main() {
    let input = include_str!("./input.txt");
    println!("{}", calc_win_possibilities(input));
}

fn calc_win_possibilities(input: &str) -> i32 {
    let mut lines = input.lines();

    let time_line = lines.next().unwrap();
    // println!("time_line = '{}'", time_line);
    let race_time = time_line
        .strip_prefix("Time:")
        .unwrap()
        .split(' ')
        .filter(|s| !s.is_empty())
        .fold(String::new(), |a, b| a + b)
        .parse::<i64>()
        .unwrap();
    // println!("race_time = {}", race_time);

    let distance_line = lines.next().unwrap();
    // println!("distance_line = '{}'", distance_line);
    let win_distance = distance_line
        .strip_prefix("Distance:")
        .unwrap()
        .split(' ')
        .filter(|s| !s.is_empty())
        .fold(String::new(), |a, b| a + b)
        .parse::<i64>()
        .unwrap();
    // println!("win_distance = {}", win_distance);

    let mut win_count = 0;
    for press_time in 1..race_time {
        let speed = press_time; // mm/ms
        let drive_time = race_time - press_time; // ms
        let distance = speed * drive_time;
        if distance > win_distance {
            win_count += 1;
        }
    }
    win_count
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_calc_win_possibilities() {
        let input = concat!("Time:      7  15   30\n", "Distance:  9  40  200\n",);
        assert_eq!(calc_win_possibilities(input), 71503);
    }
}
