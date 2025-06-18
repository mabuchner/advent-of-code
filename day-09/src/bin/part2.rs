use std::str::FromStr;

fn main() {
    let input = include_str!("./input.txt");
    println!("{}", calc_sum_of_backward_extrapolated_values(input));
}

fn calc_sum_of_backward_extrapolated_values(input: &str) -> i64 {
    let mut histories = input
        .lines()
        .map(History::from_str)
        .map(Result::unwrap)
        .collect::<Vec<_>>();
    // println!("histories = {:?}", histories);

    for history in &mut histories {
        history.extrapolate_backwards();
    }
    // println!("extrapolated histories = {:#?}", histories);

    histories
        .iter()
        .map(|h| h.values.first().unwrap())
        .sum()
}

#[derive(Debug)]
struct History {
    values: Vec<i64>,
    differences: Vec<Vec<i64>>,
}

impl FromStr for History {
    type Err = ();

    fn from_str(line: &str) -> Result<Self, Self::Err> {
        let values = line
            .split(' ')
            .map(|s| s.parse::<i64>().unwrap())
            .collect::<Vec<_>>();

        let mut differences = vec![];
        let mut base = &values;
        while !base.iter().all(|n| *n == 0) {
            let difference = base.windows(2).map(|w| w[1] - w[0]).collect::<Vec<i64>>();
            differences.push(difference);
            base = &differences[differences.len() - 1];
        }

        Ok(History {
            values,
            differences,
        })
    }
}

impl History {
    fn extrapolate_backwards(&mut self) {
        if self.differences.is_empty() {
            return;
        }

        let row_count = self.differences.len();
        self.differences.last_mut().unwrap().insert(0, 0);

        for row in (0..row_count - 1).rev() {
            let difference = *self.differences[row + 1].first().unwrap();
            let value = *self.differences[row].first().unwrap();
            self.differences[row].insert(0, value - difference);
        }

        let difference = *self.differences[0].first().unwrap();
        let value = *self.values.first().unwrap();
        self.values.insert(0, value - difference);
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_calc_sum_of_backward_extrapolated_values() {
        let input = concat!("0 3 6 9 12 15\n", "1 3 6 10 15 21\n", "10 13 16 21 30 45");
        assert_eq!(calc_sum_of_backward_extrapolated_values(input), 2);
    }
}
