fn main() {
    let input = include_str!("./input.txt");
    println!("{}", find_calibration_values_sum(input));
}

fn find_calibration_values_sum(input: &str) -> i32 {
    input.lines().map(find_calibration_value).sum::<i32>()
}

fn find_calibration_value(line: &str) -> i32 {
    let first_index = line.find(|c| char::is_ascii_digit(&c)).unwrap();
    let last_index = line.rfind(|c| char::is_ascii_digit(&c)).unwrap();

    let first_digit = line.chars().nth(first_index).unwrap();
    let last_digit = line.chars().nth(last_index).unwrap();

    let digit_str = format!("{}{}", first_digit, last_digit);
    digit_str.parse::<i32>().unwrap()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_find_calibration_values_sum() {
        let input = concat!("1abc2\n", "pqr3stu8vwx\n", "a1b2c3d4e5f\n", "treb7uchet");
        assert_eq!(find_calibration_values_sum(input), 142);
    }

    #[test]
    fn test_find_calibration_value() {
        assert_eq!(find_calibration_value("1abc2"), 12);
        assert_eq!(find_calibration_value("pqr3stu8vwx"), 38);
        assert_eq!(find_calibration_value("a1b2c3d4e5f"), 15);
        assert_eq!(find_calibration_value("treb7uchet"), 77);
    }
}
