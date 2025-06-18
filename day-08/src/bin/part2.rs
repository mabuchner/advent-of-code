use std::collections::HashMap;

fn main() {
    let input = include_str!("./input.txt");
    println!("{}", count_ghost_steps(input));
}

fn count_ghost_steps(input: &str) -> usize {
    let mut lines = input.lines().filter(|l| !l.is_empty());

    let directions = lines
        .next()
        .unwrap()
        .chars()
        .map(|c| match c {
            'L' => 0usize,
            'R' => 1usize,
            _ => panic!("unexpected direction '{}'", c),
        })
        .collect::<Vec<_>>();

    let graph = lines
        .map(|l| {
            let (from, to_str) = l.split_once(" = ").unwrap();
            let (to_left, to_right) = to_str[1..to_str.len() - 1].split_once(", ").unwrap();
            (from, to_left, to_right)
        })
        .fold(
            HashMap::<&str, [&str; 2]>::new(),
            |mut m, (from, to_left, to_right)| {
                m.insert(from, [to_left, to_right]);
                m
            },
        );
    // println!("graph = {:?}", graph);

    let starts = graph.keys().filter(|s| s.ends_with('A')).cloned();

    let end_steps = starts.map(|mut node| {
        let mut step = 0usize;
        while !node.ends_with('Z') {
            let dir_index = step % directions.len();
            let direction = directions[dir_index];
            let destinations = graph.get(node).unwrap();
            node = destinations[direction];
            step += 1;
        }
        step
    });

    end_steps.reduce(lcm).unwrap()

    /*
    // Brute force -> too slow
    let mut current = graph
        .keys()
        .filter(|s| s.ends_with('A'))
        .copied()
        .collect::<Vec<&str>>();
    println!("start = {:?}", current);
    let mut step = 0;
    while !current.iter().all(|s| s.ends_with('Z')) {
        // println!("current = {:?}", current);
        if step % 1000000 == 0 {
            println!("{}: current = {:?}", step, current);
        }
        let direction = directions[step % directions.len()];
        step += 1;

        current = current
            .iter()
            .map(|&s| {
                let next_locations = graph.get(s).unwrap();
                next_locations[direction]
            })
            .collect::<Vec<&str>>();
    }
    step
    */
}

fn lcm(a: usize, b: usize) -> usize {
    a * (b / gcd(a, b))
}

fn gcd(a: usize, b: usize) -> usize {
    if b == 0 {
        return a;
    }
    gcd(b, a % b)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_count_ghost_steps() {
        let input = concat!(
            "LR\n",
            "\n",
            "11A = (11B, XXX)\n",
            "11B = (XXX, 11Z)\n",
            "11Z = (11B, XXX)\n",
            "22A = (22B, XXX)\n",
            "22B = (22C, 22C)\n",
            "22C = (22Z, 22Z)\n",
            "22Z = (22B, 22B)\n",
            "XXX = (XXX, XXX)"
        );
        assert_eq!(count_ghost_steps(input), 6);
    }
}
