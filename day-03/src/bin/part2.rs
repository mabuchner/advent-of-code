use std::str::FromStr;

fn main() {
    let input = include_str!("./input.txt");
    println!("{}", sum_of_gear_ratios(input));
}

fn sum_of_gear_ratios(s: &str) -> i32 {
    let schematic = Schematic::from_str(s).unwrap();

    let mut sum = 0;
    for (row, schematic_row) in schematic.rows.iter().enumerate() {
        for (col, c) in schematic_row.iter().enumerate() {
            if *c == '*' {
                let nums = schematic.get_numbers_around(col as i32, row as i32);
                if nums.len() == 2 {
                    sum += nums[0] * nums[1];
                }
            }
        }
    }

    sum
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

    fn digit_at(&self, col: i32, row: i32) -> Option<char> {
        let c = self.at(col, row)?;
        if !c.is_ascii_digit() {
            return None;
        }
        Some(c)
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

    fn get_numbers_around(&self, col: i32, row: i32) -> Vec<i32> {
        let mut nums = Vec::<i32>::default();

        // Top row
        if let Some(top_center_num) = self.get_number_at(col, row - 1) {
            nums.push(top_center_num);
        } else {
            if let Some(top_left_num) = self.get_number_at(col - 1, row - 1) {
                nums.push(top_left_num);
            }
            if let Some(top_right_num) = self.get_number_at(col + 1, row - 1) {
                nums.push(top_right_num);
            }
        }

        // Bottom row
        if let Some(bottom_center_num) = self.get_number_at(col, row + 1) {
            nums.push(bottom_center_num);
        } else {
            if let Some(bottom_left_num) = self.get_number_at(col - 1, row + 1) {
                nums.push(bottom_left_num);
            }
            if let Some(bottom_right_num) = self.get_number_at(col + 1, row + 1) {
                nums.push(bottom_right_num);
            }
        }

        // Left
        if let Some(left_num) = self.get_number_at(col - 1, row) {
            nums.push(left_num);
        }

        // Right
        if let Some(right_num) = self.get_number_at(col + 1, row) {
            nums.push(right_num);
        }

        nums
    }

    fn get_number_at(&self, col: i32, row: i32) -> Option<i32> {
        self.digit_at(col, row)?;

        let mut l = col;
        while self.digit_at(l - 1, row).is_some() {
            l -= 1;
        }

        let mut r = col;
        while self.digit_at(r + 1, row).is_some() {
            r += 1;
        }

        let num_str: String = self.rows[row as usize][l as usize..(r + 1) as usize]
            .iter()
            .collect();
        match num_str.parse::<i32>() {
            Ok(num) => Some(num),
            Err(_) => None,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_sum_of_gear_ratios() {
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
        assert_eq!(sum_of_gear_ratios(input), 467835);
    }

    #[test]
    fn test_get_number_at() {
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
        let s = Schematic::from_str(input).unwrap();
        assert_eq!(s.get_number_at(0, 0), Some(467));
    }

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
