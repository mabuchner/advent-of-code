use core::fmt;
use std::{
    collections::{HashSet, VecDeque},
    str::FromStr,
};

fn main() {
    let input = include_str!("./input.txt");
    println!("{}", count_enclosed_tiles(input));
}

fn count_enclosed_tiles(input: &str) -> i32 {
    let map = Map::from_str(input).unwrap();
    println!("map =\n{:?}", map);

    let (start_col, start_row) = map.find_start().unwrap();
    // println!("start = {:?}", start);

    let path = map.get_path((start_col, start_row));
    let mut path_map = Map {
        rows: vec![vec!['.'; map.col_count as usize]; map.row_count as usize],
        col_count: map.col_count,
        row_count: map.row_count,
    };
    path.iter().for_each(|(col, row)| {
        path_map.rows[*row as usize][*col as usize] = map.rows[*row as usize][*col as usize]
    });

    let mut connected_west = false;
    if let Some(west_tile) = path_map.at_raw((start_col - 1, start_row)) {
        let (_, _, check_east, _) = Map::get_tile_connection(west_tile);
        connected_west = check_east;
    }
    let mut connected_north = false;
    if let Some(north_tile) = path_map.at_raw((start_col, start_row - 1)) {
        let (_, _, _, check_south) = Map::get_tile_connection(north_tile);
        connected_north = check_south;
    }
    let mut connected_east = false;
    if let Some(east_tile) = path_map.at_raw((start_col + 1, start_row)) {
        let (check_west, _, _, _) = Map::get_tile_connection(east_tile);
        connected_east = check_west;
    }
    let mut connected_south = false;
    if let Some(south_tile) = path_map.at_raw((start_col, start_row + 1)) {
        let (_, check_north, _, _) = Map::get_tile_connection(south_tile);
        connected_south = check_north;
    }
    if connected_north && connected_south {
        path_map.rows[start_row as usize][start_col as usize] = TILE_NORTH_SOUTH;
    } else if connected_west && connected_east {
        path_map.rows[start_row as usize][start_col as usize] = TILE_EAST_WEST;
    } else if connected_north && connected_east {
        path_map.rows[start_row as usize][start_col as usize] = TILE_NORTH_EAST;
    } else if connected_north && connected_west {
        path_map.rows[start_row as usize][start_col as usize] = TILE_NORTH_WEST;
    } else if connected_south && connected_east {
        path_map.rows[start_row as usize][start_col as usize] = TILE_SOUTH_EAST;
    } else if connected_south && connected_west {
        path_map.rows[start_row as usize][start_col as usize] = TILE_SOUTH_WEST;
    }

    println!("path_map =\n{:?}", path_map);

    let dx = filter_map_dx(&path_map);
    println!("dx =\n{:?}", dx);

    let dy = filter_map_dy(&path_map);
    println!("dy =\n{:?}", dy);

    let mut inside_map = Map {
        rows: vec![vec!['.'; map.col_count as usize]; map.row_count as usize],
        col_count: map.col_count,
        row_count: map.row_count,
    };

    let mut inside_count = 0;
    for row in 0..path_map.row_count as usize {
        for col in 0..path_map.col_count as usize {
            let t = path_map.rows[row][col];
            if t != '.' {
                //inside_map.rows[row][col] = t;
                continue;
            }

            // Ray west
            let mut west_in = false;
            for i in 0..col {
                let t = dx.rows[row][i];
                if t != '.' {
                    west_in = !west_in;
                }
            }
            // Ray east
            let mut east_in = false;
            for i in col + 1..path_map.col_count as usize {
                let t = dx.rows[row][i];
                if t != '.' {
                    east_in = !east_in;
                }
            }
            // Ray north
            let mut north_in = false;
            for i in 0..row {
                let t = dy.rows[i][col];
                if t != '.' {
                    north_in = !north_in;
                }
            }
            // Ray west
            let mut south_in = false;
            for i in row + 1..path_map.row_count as usize {
                let t = dy.rows[i][col];
                if t != '.' {
                    south_in = !south_in;
                }
            }

            // if west_in && east_in && north_in && south_in {
            if (west_in && north_in) {
                inside_map.rows[row][col] = 'I';
                inside_count += 1;
            }
        }
    }
    println!("inside_map =\n{:?}", inside_map);

    inside_count
}

fn filter_map_dx(map: &Map) -> Map {
    let mut dx = Map {
        rows: vec![vec!['.'; map.col_count as usize]; map.row_count as usize],
        col_count: map.col_count,
        row_count: map.row_count,
    };
    /*
    for row in 0..dx.row_count as usize {
        for col in 0..dx.col_count as usize {
            let t = map.rows[row][col];
            if t == TILE_NORTH_SOUTH
                || t == TILE_NORTH_EAST
                || t == TILE_NORTH_WEST
                || t == TILE_SOUTH_WEST
                || t == TILE_SOUTH_EAST
            {
                dx.rows[row][col] = TILE_NORTH_SOUTH;
            }
        }
    }
    */
    for row in 0..dx.row_count as usize {
        let mut last_tile = '.';
        for col in 0..dx.col_count as usize {
            let t = map.rows[row][col];
            let next_tile = if col + 1 < dx.col_count as usize {
                map.rows[row][col + 1]
            } else {
                '.'
            };
            if t != '.' && (!is_connected_hor(last_tile, t) || !is_connected_hor(t, next_tile)) {
                dx.rows[row][col] = TILE_NORTH_SOUTH;
            }
            last_tile = t;
        }
    }
    dx
}

fn filter_map_dy(map: &Map) -> Map {
    let mut dy = Map {
        rows: vec![vec!['.'; map.col_count as usize]; map.row_count as usize],
        col_count: map.col_count,
        row_count: map.row_count,
    };
    /*
    for row in 0..dy.row_count as usize {
        for col in 0..dy.col_count as usize {
            let t = map.rows[row][col];
            if t == TILE_EAST_WEST
                || t == TILE_NORTH_EAST
                || t == TILE_NORTH_WEST
                || t == TILE_SOUTH_WEST
                || t == TILE_SOUTH_EAST
                || t == TILE_SOUTH_EAST
            {
                dy.rows[row][col] = TILE_EAST_WEST;
            }
        }
    }
    */
    for col in 0..dy.col_count as usize {
        let mut last_tile = '.';
        for row in 0..dy.row_count as usize {
            let t = map.rows[row][col];
            let next_tile = if row + 1 < dy.row_count as usize {
                map.rows[row + 1][col]
            } else {
                '.'
            };
            if t != '.' && (!is_connected_vert(last_tile, t) || !is_connected_vert(t, next_tile)) {
                dy.rows[row][col] = TILE_EAST_WEST;
            }
            last_tile = t;
        }
    }
    dy
}

fn is_connected_hor(left: char, right: char) -> bool {
    if left == TILE_NORTH_SOUTH
        || right == TILE_NORTH_SOUTH
        || left == TILE_NORTH_WEST
        || left == TILE_SOUTH_WEST
        || left == '.'
        || right == '.'
    {
        return false;
    }
    if left == TILE_EAST_WEST {
        return right == TILE_EAST_WEST || right == TILE_NORTH_WEST || right == TILE_SOUTH_WEST;
    }
    if left == TILE_NORTH_EAST {
        return right == TILE_EAST_WEST || right == TILE_NORTH_WEST || right == TILE_SOUTH_WEST;
    }
    if left == TILE_SOUTH_EAST {
        return right == TILE_EAST_WEST || right == TILE_NORTH_WEST || right == TILE_SOUTH_WEST;
    }
    false
}

fn is_connected_vert(top: char, bottom: char) -> bool {
    if top == TILE_EAST_WEST
        || bottom == TILE_EAST_WEST
        || top == TILE_NORTH_EAST
        || top == TILE_NORTH_WEST
        || top == '.'
        || bottom == '.'
    {
        return false;
    }
    if top == TILE_NORTH_SOUTH {
        return bottom == TILE_NORTH_SOUTH
            || bottom == TILE_NORTH_WEST
            || bottom == TILE_NORTH_EAST;
    }
    if top == TILE_SOUTH_EAST {
        return bottom == TILE_NORTH_SOUTH
            || bottom == TILE_NORTH_EAST
            || bottom == TILE_NORTH_WEST;
    }
    if top == TILE_SOUTH_WEST {
        return bottom == TILE_NORTH_SOUTH
            || bottom == TILE_NORTH_EAST
            || bottom == TILE_NORTH_WEST;
    }
    false
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

    fn get_path(&self, (start_col, start_row): (i32, i32)) -> Vec<(i32, i32)> {
        let mut visited = HashSet::<(i32, i32)>::new();
        let mut stack = VecDeque::from([(start_col, start_row, i32::MIN, i32::MIN, vec![])]);
        let mut final_path = vec![];
        while let Some((col, row, last_col, last_row, path)) = stack.pop_back() {
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
                    final_path = path;
                    //println!("max_loop_length = {}", max_loop_length);
                    break;
                }
                continue;
            }
            visited.insert((col, row));

            let mut new_path = path.clone();
            new_path.push((col, row));

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
                            stack.push_back((new_col, new_row, col, row, new_path.clone()));
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
                            stack.push_back((new_col, new_row, col, row, new_path.clone()));
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
                            stack.push_back((new_col, new_row, col, row, new_path.clone()));
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
                            stack.push_back((new_col, new_row, col, row, new_path.clone()));
                        }
                    }
                }
            }
        }
        final_path
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_count_enclosed_tiles() {
        let input0 = concat!(
            "...........\n",
            ".S-------7.\n",
            ".|F-----7|.\n",
            ".||.....||.\n",
            ".||.....||.\n",
            ".|L-7.F-J|.\n",
            ".|..|.|..|.\n",
            ".L--J.L--J.\n",
            "..........."
        );
        let input1 = concat!(
            ".F----7F7F7F7F-7....\n",
            ".|F--7||||||||FJ....\n",
            ".||.FJ||||||||L7....\n",
            "FJL7L7LJLJ||LJ.L-7..\n",
            "L--J.L7...LJS7F-7L7.\n",
            "....F-J..F7FJ|L7L7L7\n",
            "....L7.F7||L7|.L7L7|\n",
            ".....|FJLJ|FJ|F7|.LJ\n",
            "....FJL-7.||.||||...\n",
            "....L---J.LJ.LJLJ..."
        );
        let input2 = concat!(
            "FF7FSF7F7F7F7F7F---7\n",
            "L|LJ||||||||||||F--J\n",
            "FL-7LJLJ||||||LJL-77\n",
            "F--JF--7||LJLJ7F7FJ-\n",
            "L---JF-JLJ.||-FJLJJ7\n",
            "|F|F-JF---7F7-L7L|7|\n",
            "|FFJF7L7F-JF7|JL---7\n",
            "7-L-JL7||F7|L7F-7F7|\n",
            "L.L7LFJ|||||FJL7||LJ\n",
            "L7JLJL-JLJLJL--JLJ.L"
        );
        assert_eq!(count_enclosed_tiles(input0), 4);
        assert_eq!(count_enclosed_tiles(input1), 8);
        assert_eq!(count_enclosed_tiles(input2), 10);
    }
}
