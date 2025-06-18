use std::{cmp::Ordering, collections::HashMap};

#[macro_use]
extern crate lazy_static;

fn main() {
    let input = include_str!("./input.txt");
    println!("{}", calc_total_winnings(input));
}

fn calc_total_winnings(input: &str) -> i32 {
    let mut hands_and_bets = input
        .lines()
        .map(|s| {
            s.split_once(' ')
                .map(|(hand, bet_str)| (hand, bet_str.parse::<i32>().unwrap(), get_card_kind(hand)))
                .unwrap()
        })
        .collect::<Vec<_>>();
    // println!("hands_and_bets = {:?}", hands_and_bets);

    hands_and_bets.sort_by(|(hand_a, _bet_a, kind_a), (hand_b, _bet_b, kind_b)| {
        if kind_a < kind_b {
            return Ordering::Less;
        }

        if kind_a > kind_b {
            return Ordering::Greater;
        }

        for (ca, cb) in hand_a.chars().zip(hand_b.chars()) {
            let strength_a = CARD_STRENGTH.get(&ca);
            let strength_b = CARD_STRENGTH.get(&cb);
            if strength_a < strength_b {
                return Ordering::Less;
            }
            if strength_a > strength_b {
                return Ordering::Greater;
            }
        }

        Ordering::Equal
    });
    // println!("sorted hands_and_bets = {:?}", hands_and_bets);

    hands_and_bets
        .iter()
        .zip(1..)
        .fold(0, |sum, ((_hand, bet, _kind), rank)| sum + rank * bet)
}

lazy_static! {
    static ref CARD_STRENGTH: HashMap<char, i32> = {
        let mut m = HashMap::new();
        m.insert('2', 1);
        m.insert('3', 2);
        m.insert('4', 3);
        m.insert('5', 4);
        m.insert('6', 5);
        m.insert('7', 6);
        m.insert('8', 7);
        m.insert('9', 8);
        m.insert('T', 9);
        m.insert('J', 10);
        m.insert('Q', 11);
        m.insert('K', 12);
        m.insert('A', 13);
        m
    };
}

#[derive(Debug, PartialEq, Eq, PartialOrd)]
enum CardKind {
    HighCard = 1,
    OnePair,
    TwoPair,
    ThreeOfAKind,
    FullHouse,
    FourOfAKind,
    FiveOfAKind,
}

fn get_card_kind(hand: &str) -> CardKind {
    let counts = hand.chars().fold(HashMap::<char, i32>::new(), |mut m, c| {
        m.entry(c).and_modify(|counter| *counter += 1).or_insert(1);
        m
    });

    // AAAAA
    if counts.len() == 1 {
        return CardKind::FiveOfAKind;
    }

    // ABCDE
    if counts.len() == 5 {
        return CardKind::HighCard;
    }

    // AABCD
    if counts.len() == 4 {
        return CardKind::OnePair;
    }

    // AAABC || AABBC
    if counts.len() == 3 {
        let mut values = counts.values();
        if *values.next().unwrap() == 3
            || *values.next().unwrap() == 3
            || *values.next().unwrap() == 3
        {
            return CardKind::ThreeOfAKind;
        }
        return CardKind::TwoPair;
    }

    // AAAAB || AABBB
    if counts.len() == 2 {
        let mut values = counts.values();
        if *values.next().unwrap() == 4 || *values.next().unwrap() == 4 {
            return CardKind::FourOfAKind;
        }
        return CardKind::FullHouse;
    }

    CardKind::HighCard
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_total_winnings() {
        let input = concat!(
            "32T3K 765\n",
            "T55J5 684\n",
            "KK677 28\n",
            "KTJJT 220\n",
            "QQQJA 483\n",
        );
        assert_eq!(calc_total_winnings(input), 6440);
    }
}
