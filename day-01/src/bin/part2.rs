fn main() {
    let input = include_str!("./input.txt");
    println!("{}", find_calibration_values_sum(input));
}

struct WordAndDigit<'a> {
    word: &'a str,
    digit: char,
}

const NUMBER_WORDS: &[WordAndDigit] = &[
    WordAndDigit {
        word: "1",
        digit: '1',
    },
    WordAndDigit {
        word: "2",
        digit: '2',
    },
    WordAndDigit {
        word: "3",
        digit: '3',
    },
    WordAndDigit {
        word: "4",
        digit: '4',
    },
    WordAndDigit {
        word: "5",
        digit: '5',
    },
    WordAndDigit {
        word: "6",
        digit: '6',
    },
    WordAndDigit {
        word: "7",
        digit: '7',
    },
    WordAndDigit {
        word: "8",
        digit: '8',
    },
    WordAndDigit {
        word: "9",
        digit: '9',
    },
    WordAndDigit {
        word: "one",
        digit: '1',
    },
    WordAndDigit {
        word: "two",
        digit: '2',
    },
    WordAndDigit {
        word: "three",
        digit: '3',
    },
    WordAndDigit {
        word: "four",
        digit: '4',
    },
    WordAndDigit {
        word: "five",
        digit: '5',
    },
    WordAndDigit {
        word: "six",
        digit: '6',
    },
    WordAndDigit {
        word: "seven",
        digit: '7',
    },
    WordAndDigit {
        word: "eight",
        digit: '8',
    },
    WordAndDigit {
        word: "nine",
        digit: '9',
    },
];

fn find_calibration_values_sum(lines: &str) -> i32 {
    lines.lines().map(find_calibration_value).sum::<i32>()
}

fn find_calibration_value(line: &str) -> i32 {
    let first_digit = find_first_digit(line).unwrap();
    let last_digit = find_last_digit(line).unwrap();
    let digit_str = format!("{}{}", first_digit, last_digit);
    digit_str.parse::<i32>().unwrap()
}

fn find_first_digit(line: &str) -> Option<char> {
    NUMBER_WORDS
        .iter()
        .map(|word_and_digit| (line.find(word_and_digit.word), word_and_digit.digit))
        .filter(|(i, _)| i.is_some())
        .map(|(i, d)| (i.unwrap(), d))
        .min_by(|(ia, _), (ib, _)| ia.cmp(ib))
        .map(|(_, d)| d)
}

fn find_last_digit(line: &str) -> Option<char> {
    NUMBER_WORDS
        .iter()
        .map(|word_and_digit| (line.rfind(word_and_digit.word), word_and_digit.digit))
        .filter(|(i, _)| i.is_some())
        .map(|(i, d)| (i.unwrap(), d))
        .max_by(|(ia, _), (ib, _)| ia.cmp(ib))
        .map(|(_, d)| d)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_find_calibration_values_sum() {
        let lines = concat!(
            "two1nine\n",
            "eightwothree\n",
            "abcone2threexyz\n",
            "xtwone3four\n",
            "4nineeightseven2\n",
            "zoneight234\n",
            "7pqrstsixteen"
        );
        assert_eq!(find_calibration_values_sum(lines), 281);
    }

    #[test]
    fn test_find_calibration_value() {
        assert_eq!(find_calibration_value("two1nine"), 29);
        assert_eq!(find_calibration_value("eightwothree"), 83);
        assert_eq!(find_calibration_value("abcone2threexyz"), 13);
        assert_eq!(find_calibration_value("xtwone3four"), 24);
        assert_eq!(find_calibration_value("4nineeightseven2"), 42);
        assert_eq!(find_calibration_value("zoneight234"), 14);
        assert_eq!(find_calibration_value("7pqrstsixteen"), 76);
    }

    #[test]
    fn test_find_first_digit() {
        assert_eq!(find_first_digit(""), None);
        assert_eq!(find_first_digit("two1nine"), Some('2'));
        assert_eq!(find_first_digit("eightwothree"), Some('8'));
        assert_eq!(find_first_digit("abcone2threexyz"), Some('1'));
        assert_eq!(find_first_digit("xtwone3four"), Some('2'));
        assert_eq!(find_first_digit("4nineeightseven2"), Some('4'));
        assert_eq!(find_first_digit("zoneight234"), Some('1'));
        assert_eq!(find_first_digit("7pqrstsixteen"), Some('7'));
    }

    #[test]
    fn test_find_last_digit() {
        assert_eq!(find_last_digit(""), None);
        assert_eq!(find_last_digit("two1nine"), Some('9'));
        assert_eq!(find_last_digit("eightwothree"), Some('3'));
        assert_eq!(find_last_digit("abcone2threexyz"), Some('3'));
        assert_eq!(find_last_digit("xtwone3four"), Some('4'));
        assert_eq!(find_last_digit("4nineeightseven2"), Some('2'));
        assert_eq!(find_last_digit("zoneight234"), Some('4'));
        assert_eq!(find_last_digit("7pqrstsixteen"), Some('6'));
    }
}
