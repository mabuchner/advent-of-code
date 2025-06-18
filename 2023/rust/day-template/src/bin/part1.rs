fn main() {
    let input = include_str!("./input.txt");
    println!("{}", count_steps_to_zzz(input));
}

fn count_steps_to_zzz(input: &str) -> i32 {
    0
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_count_steps_to_zzz() {
        let input = concat!(
            "RL\n",
            "\n",
            "AAA = (BBB, CCC)\n",
            "BBB = (DDD, EEE)\n",
            "CCC = (ZZZ, GGG)\n",
            "DDD = (DDD, DDD)\n",
            "EEE = (EEE, EEE)\n",
            "GGG = (GGG, GGG)\n",
            "ZZZ = (ZZZ, ZZZ)\n",
        );
        assert_eq!(count_steps_to_zzz(input), 142);
    }
}
