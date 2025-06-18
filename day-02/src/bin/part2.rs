use std::{cmp, str::FromStr};

fn main() {
    let input = include_str!("./input.txt");
    println!("{}", find_min_cubes_power_sum(input));
}

fn find_min_cubes_power_sum(input: &str) -> i32 {
    input
        .lines()
        .map(Game::from_str)
        .map(|g| g.unwrap())
        .map(|g| g.max_grab())
        .map(|g| g.power())
        .sum()
}

fn add_possible_game_ids(input: &str, max_red: i32, max_green: i32, max_blue: i32) -> i32 {
    input
        .lines()
        .map(Game::from_str)
        .map(|g| g.unwrap())
        .filter(|g| g.is_possible(max_red, max_green, max_blue))
        .fold(0, |id_sum, gb| id_sum + gb.id)
}

struct Game {
    id: i32,
    grabs: Vec<Grab>,
}

impl Game {
    fn is_possible(&self, max_red: i32, max_green: i32, max_blue: i32) -> bool {
        self.grabs
            .iter()
            .all(|g| g.is_possible(max_red, max_green, max_blue))
    }

    fn max_grab(&self) -> Grab {
        self.grabs.iter().fold(
            Grab {
                red: 0,
                green: 0,
                blue: 0,
            },
            |a, b| Grab {
                red: cmp::max(a.red, b.red),
                green: cmp::max(a.green, b.green),
                blue: cmp::max(a.blue, b.blue),
            },
        )
    }
}

impl FromStr for Game {
    type Err = std::num::ParseIntError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let colon_idx = s.find(": ").unwrap();

        let id_str = &s["Game ".len()..colon_idx];
        let id = id_str.parse::<i32>()?;

        //Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue
        let grabs_str = &s[colon_idx + ": ".len()..];
        let grabs: Vec<Grab> = grabs_str
            .split("; ")
            .map(|s| Grab::from_str(s).unwrap())
            .collect();

        Ok(Self { id, grabs })
    }
}

#[derive(Debug, PartialEq)]
struct Grab {
    red: i32,
    green: i32,
    blue: i32,
}

impl Grab {
    fn is_possible(&self, max_red: i32, max_green: i32, max_blue: i32) -> bool {
        self.red <= max_red && self.green <= max_green && self.blue <= max_blue
    }

    fn power(&self) -> i32 {
        self.red * self.green * self.blue
    }
}

impl FromStr for Grab {
    type Err = std::num::ParseIntError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let (red, green, blue) = s
            .split(", ")
            .map(|color_str| {
                if let Some(num_str) = color_str.strip_suffix(" red") {
                    return (num_str.parse::<i32>().unwrap(), 0, 0);
                }
                if let Some(num_str) = color_str.strip_suffix(" green") {
                    return (0, num_str.parse::<i32>().unwrap(), 0);
                }
                if let Some(num_str) = color_str.strip_suffix(" blue") {
                    return (0, 0, num_str.parse::<i32>().unwrap());
                }
                (0, 0, 0)
            })
            .reduce(|(red_a, green_a, blue_a), (red_b, green_b, blue_b)| {
                (red_a + red_b, green_a + green_b, blue_a + blue_b)
            })
            .unwrap();

        Ok(Grab { red, green, blue })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_add_possible_game_ids() {
        let lines = concat!(
            "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green\n",
            "Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue\n",
            "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red\n",
            "Game 4: 1 green, 3 red, 6 blue; 3 green, 6 red; 3 green, 15 blue, 14 red\n",
            "Game 5: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green"
        );
        assert_eq!(add_possible_game_ids(lines, 12, 13, 14), 8);
    }

    #[test]
    fn test_game_from_str() {
        let g1 = Game::from_str("Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green").unwrap();
        assert_eq!(g1.id, 1);
        assert_eq!(g1.grabs.len(), 3);
    }

    #[test]
    fn test_game_is_possible() {
        assert!(Game {
            id: 1,
            grabs: vec![
                Grab {
                    red: 4,
                    green: 0,
                    blue: 3,
                },
                Grab {
                    red: 1,
                    green: 2,
                    blue: 6,
                },
                Grab {
                    red: 0,
                    green: 2,
                    blue: 0,
                },
            ],
        }
        .is_possible(12, 13, 14));

        assert!(!Game {
            id: 1,
            grabs: vec![
                Grab {
                    red: 20,
                    green: 8,
                    blue: 6,
                },
                Grab {
                    red: 4,
                    green: 13,
                    blue: 5,
                },
                Grab {
                    red: 1,
                    green: 5,
                    blue: 0,
                },
            ],
        }
        .is_possible(12, 13, 14));
    }

    #[test]
    fn test_game_max_grab() {
        assert_eq!(
            Game {
                id: 1,
                grabs: vec![
                    Grab {
                        red: 4,
                        green: 0,
                        blue: 3,
                    },
                    Grab {
                        red: 1,
                        green: 2,
                        blue: 6,
                    },
                    Grab {
                        red: 0,
                        green: 2,
                        blue: 0,
                    },
                ],
            }
            .max_grab(),
            Grab {
                red: 4,
                green: 2,
                blue: 6
            }
        );
    }

    #[test]
    fn test_grab_from_str() {
        let g1 = Grab::from_str("1 red, 2 green, 6 blue").unwrap();
        assert_eq!(
            g1,
            Grab {
                red: 1,
                green: 2,
                blue: 6
            }
        );
    }

    #[test]
    fn test_grab_is_possible() {
        assert!(Grab {
            red: 1,
            green: 2,
            blue: 6
        }
        .is_possible(12, 13, 14));

        assert!(!Grab {
            red: 20,
            green: 8,
            blue: 6
        }
        .is_possible(12, 13, 14));
    }

    #[test]
    fn test_grab_power() {
        assert_eq!(
            Grab {
                red: 20,
                green: 8,
                blue: 6
            }
            .power(),
            960
        );
    }
}
