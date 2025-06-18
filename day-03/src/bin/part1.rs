use std::str::FromStr;

fn main() {
    let input = include_str!("./input.txt");
    println!("{}", sum_of_engine_parts(input));
}

fn sum_of_engine_parts(s: &str) -> i32 {
    let schematic = Schematic::from_str(s).unwrap();

    let mut sum = 0;
    for (row, schematic_row) in schematic.rows.iter().enumerate() {
        let mut num_str = String::default();
        let mut symbol_around = false;
        for (col, c) in schematic_row.iter().enumerate() {
            if c.is_ascii_digit() {
                num_str.push(*c);
                if schematic.is_symbol_around(col as i32, row as i32) {
                    symbol_around = true;
                }
            } else {
                if !num_str.is_empty() && symbol_around {
                    sum += num_str.parse::<i32>().unwrap();
                }
                num_str.clear();
                symbol_around = false;
            }
        }
        if !num_str.is_empty() && symbol_around {
            sum += num_str.parse::<i32>().unwrap();
        }
    }
    sum
}

type SchematicRow = Vec<char>;

struct Schematic {
    col_count: i32,
    row_count: i32,
    rows: Vec<SchematicRow>,
}

impl FromStr for Schematic {
    type Err = ();

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let rows: Vec<SchematicRow> = s
            .lines()
            .map(|l| l.chars().collect::<SchematicRow>())
            .collect();

        let col_count = i32::try_from(rows.first().unwrap().len()).unwrap();
        let row_count = i32::try_from(rows.len()).unwrap();

        Ok(Self {
            col_count,
            row_count,
            rows,
        })
    }
}

impl Schematic {
    fn at(&self, col: i32, row: i32) -> Option<char> {
        if col < 0 || col >= self.col_count || row < 0 || row >= self.row_count {
            return None;
        }
        Some(self.rows[row as usize][col as usize])
    }

    fn is_symbol_at(&self, col: i32, row: i32) -> bool {
        let c = self.at(col, row).unwrap_or('.');
        c != '.' && !c.is_alphanumeric()
    }

    fn is_symbol_around(&self, col: i32, row: i32) -> bool {
        self.is_symbol_at(col - 1, row)
            || self.is_symbol_at(col - 1, row - 1)
            || self.is_symbol_at(col, row - 1)
            || self.is_symbol_at(col + 1, row - 1)
            || self.is_symbol_at(col + 1, row)
            || self.is_symbol_at(col + 1, row + 1)
            || self.is_symbol_at(col, row + 1)
            || self.is_symbol_at(col - 1, row + 1)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_sum_of_engine_parts() {
        let input = concat!(
            "467..114..\n",
            "...*......\n",
            "..35..633.\n",
            "......#...\n",
            "617*......\n",
            ".....+.58.\n",
            "..592.....\n",
            "......755.\n",
            "...$.*....\n",
            ".664.598..\n",
        );
        assert_eq!(sum_of_engine_parts(input), 4361);
    }
}
