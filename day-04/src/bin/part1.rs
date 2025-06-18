use std::collections::HashSet;

use lazy_static::lazy_static;
use regex::Regex;

fn main() {
    let input = include_str!("./input.txt");
    println!("{}", calc_card_points_sum(input));
}

lazy_static! {
    static ref RE: Regex = Regex::new(r"Card\s+\d+?: (.+) \| (.+)").unwrap();
}

fn calc_card_points_sum(input: &str) -> i32 {
    RE.captures_iter(input)
        .map(|c| c.extract())
        .map(|(_, [owning_str, winning_str])| {
            // println!("'{}', '{}'", owning_str, winning_str);
            let owning_nums: HashSet<i32> = owning_str
                .split(' ')
                .filter(|s| !s.is_empty())
                .map(|s| s.parse::<i32>().unwrap())
                .collect();
            let winning_nums: HashSet<i32> = winning_str
                .split(' ')
                .filter(|s| !s.is_empty())
                .map(|s| s.parse::<i32>().unwrap())
                .collect();
            // println!("{:?}, {:?}", owning_nums, winning_nums);

            let intersection = owning_nums.intersection(&winning_nums);
            // println!("i = {:?}", intersection);
            let intersection_count = intersection.count();

            let mut points = 0;
            if intersection_count > 0 {
                points = 1 << (intersection_count - 1);
            }
            points
        })
        .sum()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_calc_card_points_sum() {
        let input = concat!(
            "Card 1: 41 48 83 86 17 | 83 86  6 31 17  9 48 53\n",
            "Card 2: 13 32 20 16 61 | 61 30 68 82 17 32 24 19\n",
            "Card 3:  1 21 53 59 44 | 69 82 63 72 16 21 14  1\n",
            "Card 4: 41 92 73 84 69 | 59 84 76 51 58  5 54 83\n",
            "Card 5: 87 83 26 28 32 | 88 30 70 12 93 22 82 36\n",
            "Card 6: 31 18 13 56 72 | 74 77 10 23 35 67 36 11\n"
        );
        assert_eq!(calc_card_points_sum(input), 13);
    }
}
