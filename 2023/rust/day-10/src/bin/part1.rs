use core::fmt;
use std::{
    collections::{HashSet, VecDeque},
    str::FromStr,
};

fn main() {
    let input = include_str!("./input.txt");
    println!("{}", find_furthest_distance(input));
}

fn find_furthest_distance(input: &str) -> i32 {
    let map = Map::from_str(input).unwrap();
    println!("map =\n{:?}", map);

    let start = map.find_start().unwrap();
    // println!("start = {:?}", start);

    map.explore(start)
}

const TILE_NORTH_SOUTH: char = '|';
const TILE_EAST_WEST: char = '-';
const TILE_NORTH_EAST: char = 'L';
const TILE_NORTH_WEST: char = 'J';
const TILE_SOUTH_WEST: char = '7';
const TILE_SOUTH_EAST: char = 'F';
const TILE_START: char = 'S';

type Row = Vec<char>;

struct Map {
    rows: Vec<Row>,
    col_count: i32,
    row_count: i32,
}

impl FromStr for Map {
    type Err = ();

    fn from_str(input: &str) -> Result<Self, Self::Err> {
        let rows = input
            .lines()
            .map(|l| l.chars().collect::<Vec<_>>())
            .collect::<Vec<_>>();
        let row_count = rows.len() as i32;
        let col_count = rows[0].len() as i32;
        Ok(Map {
            rows,
            col_count,
            row_count,
        })
    }
}

impl fmt::Debug for Map {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        self.rows.iter().for_each(|row| {
            row.iter().for_each(|c| {
                write!(f, "{}", c).unwrap();
            });
            writeln!(f).unwrap();
        });
        Ok(())
    }
}

impl Map {
    fn at_raw(&self, (col, row): (i32, i32)) -> Option<char> {
        if col < 0 || col >= self.col_count || row < 0 || row >= self.row_count {
            return None;
        }
        Some(self.rows[row as usize][col as usize])
    }

    fn find_start(&self) -> Option<(i32, i32)> {
        for row in 0..self.row_count {
            for col in 0..self.col_count {
                if let Some(c) = self.at_raw((col, row)) {
                    if c == TILE_START {
                        return Some((col, row));
                    }
                }
            }
        }
        None
    }

    fn get_tile_connection(t: char) -> (bool, bool, bool, bool) {
        match t {
            TILE_NORTH_SOUTH => (false, true, false, true),
            TILE_EAST_WEST => (true, false, true, false),
            TILE_NORTH_EAST => (false, true, true, false),
            TILE_NORTH_WEST => (true, true, false, false),
            TILE_SOUTH_EAST => (false, false, true, true),
            TILE_SOUTH_WEST => (true, false, false, true),
            TILE_START => (true, true, true, true),
            _ => (false, false, false, false),
        }
    }

    fn explore(&self, (start_col, start_row): (i32, i32)) -> i32 {
        let mut visited = HashSet::<(i32, i32)>::new();
        let mut stack = VecDeque::from([(start_col, start_row, i32::MIN, i32::MIN, 0)]);
        let mut max_loop_length = 0;
        while let Some((col, row, last_col, last_row, distance)) = stack.pop_back() {
            //println!("checking {}, {}, {}", col, row, distance);
            let t = self.at_raw((col, row));
            if t.is_none() {
                continue;
            }
            let tt = t.unwrap();

            if visited.contains(&(col, row)) {
                //println!("visited");
                if tt == TILE_START {
                    //println!("distance = {}", distance);
                    max_loop_length = distance;
                    //println!("max_loop_length = {}", max_loop_length);
                    break;
                }
                continue;
            }
            visited.insert((col, row));

            let (check_west, check_north, check_east, check_south) = Map::get_tile_connection(tt);

            if check_west {
                let new_col = col - 1;
                let new_row = row;
                if new_col != last_col || new_row != last_row {
                    if let Some(west_tile) = self.at_raw((new_col, new_row)) {
                        if west_tile == TILE_START
                            || west_tile == TILE_EAST_WEST
                            || west_tile == TILE_NORTH_EAST
                            || west_tile == TILE_SOUTH_EAST
                        {
                            stack.push_back((new_col, new_row, col, row, distance + 1));
                        }
                    }
                }
            }
            if check_north {
                let new_col = col;
                let new_row = row - 1;
                if new_col != last_col || new_row != last_row {
                    if let Some(north_tile) = self.at_raw((new_col, new_row)) {
                        if north_tile == TILE_START
                            || north_tile == TILE_NORTH_SOUTH
                            || north_tile == TILE_SOUTH_WEST
                            || north_tile == TILE_SOUTH_EAST
                        {
                            stack.push_back((new_col, new_row, col, row, distance + 1));
                        }
                    }
                }
            }
            if check_east {
                let new_col = col + 1;
                let new_row = row;
                if new_col != last_col || new_row != last_row {
                    if let Some(east_tile) = self.at_raw((new_col, new_row)) {
                        if east_tile == TILE_START
                            || east_tile == TILE_EAST_WEST
                            || east_tile == TILE_NORTH_WEST
                            || east_tile == TILE_SOUTH_WEST
                        {
                            stack.push_back((new_col, new_row, col, row, distance + 1));
                        }
                    }
                }
            }
            if check_south {
                let new_col = col;
                let new_row = row + 1;
                if new_col != last_col || new_row != last_row {
                    if let Some(south_tile) = self.at_raw((new_col, new_row)) {
                        if south_tile == TILE_START
                            || south_tile == TILE_NORTH_SOUTH
                            || south_tile == TILE_NORTH_EAST
                            || south_tile == TILE_NORTH_WEST
                        {
                            stack.push_back((new_col, new_row, col, row, distance + 1));
                        }
                    }
                }
            }
        }
        // println!("max_loop = {}", max_loop_length);
        (max_loop_length + 1) / 2
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_find_furthest_distance() {
        let input0_a = concat!(".....\n", ".S-7.\n", ".|.|.\n", ".L-J.\n", ".....",);
        let input0_b = concat!("-L|F7\n", "7S-7|\n", "L|7||\n", "-L-J|\n", "L|-JF");
        let input2 = concat!("..F7.\n", ".FJ|.\n", "SJ.L7\n", "|F--J\n", "LJ...");
        assert_eq!(find_furthest_distance(input0_a), 4);
        assert_eq!(find_furthest_distance(input0_b), 4);
        assert_eq!(find_furthest_distance(input2), 8);
    }
}
