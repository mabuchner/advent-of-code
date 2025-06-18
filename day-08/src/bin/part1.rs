use std::collections::HashMap;

fn main() {
    let input = include_str!("./input.txt");
    println!("{}", count_steps_from_aaa_to_zzz(input));
}

fn count_steps_from_aaa_to_zzz(input: &str) -> usize {
    let mut lines = input.lines().filter(|l| !l.is_empty());

    let directions = lines.next().unwrap().chars().collect::<Vec<_>>();

    let graph = lines
        .map(|l| {
            let (from, to_str) = l.split_once(" = ").unwrap();
            let (to_left, to_right) = to_str[1..to_str.len() - 1].split_once(", ").unwrap();
            (from, to_left, to_right)
        })
        .fold(
            HashMap::<&str, (&str, &str)>::new(),
            |mut m, (from, to_left, to_right)| {
                m.insert(from, (to_left, to_right));
                m
            },
        );
    // println!("graph = {:?}", graph);

    let mut step = 0;
    let mut current = "AAA";
    while current != "ZZZ" {
        let direction = directions[step % directions.len()];
        step += 1;
        match direction {
            'L' => (current, _) = graph[current],
            'R' => (_, current) = graph[current],
            _ => panic!("unexpected direction '{}'", direction),
        }
    }
    step
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_count_steps_from_aaa_to_zzz() {
        let input0 = concat!(
            "RL\n",
            "\n",
            "AAA = (BBB, CCC)\n",
            "BBB = (DDD, EEE)\n",
            "CCC = (ZZZ, GGG)\n",
            "DDD = (DDD, DDD)\n",
            "EEE = (EEE, EEE)\n",
            "GGG = (GGG, GGG)\n",
            "ZZZ = (ZZZ, ZZZ)",
        );
        let input1 = concat!(
            "LLR\n",
            "\n",
            "AAA = (BBB, BBB)\n",
            "BBB = (AAA, ZZZ)\n",
            "ZZZ = (ZZZ, ZZZ)"
        );
        assert_eq!(count_steps_from_aaa_to_zzz(input0), 2);
        assert_eq!(count_steps_from_aaa_to_zzz(input1), 6);
    }
}
