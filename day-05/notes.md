seed_ranges = [Range { start: 79, len: 14 }, Range { start: 55, len: 13 }]
soil_ranges = [Range { start: 81, len: 14 }, Range { start: 57, len: 13 }]
fertilizer_ranges = [Range { start: 81, len: 14 }, Range { start: 57, len: 13 }]

seeds:
79 14
55 13

seed-to-soil map:
50 98 2
52 50 48

seed-to-soil-ranges
a 98 - 99 -> 50 - 51
b 50 - 97 -> 52 - 99

seed_ranges
79 - 92
-> no intersection with a -> 79 - 92
-> intersection with b -> 81 - 94 ok

55 - 67 ->
-> no intersection with a -> 55 - 67
-> intersection with b -> 57 - 69 ok


soil-to-fertilizer map:
a 0 15 37 -> 15 - 51
b 37 52 2 -> 52 - 53
c 39 0 15 -> 0 - 14

81 - 94
-> no intersection with a
-> no intersection with b
-> no intersection with c

57 - 69
-> no intersection with a
-> no intersection with b
-> no intersection with c

fertilizer-to-water map:
a 49 53 8 -> 53 - 60 -> 49 - 56
b 0 11 42 -> 11 - 52 -> 0 - 41
c 42 0 7 -> 0 - 6 -> 42 - 48
d 57 7 4 -> 7 - 10 -> 57 - 60


81 - 94
-> no intersection with a
-> no intersection with b
-> no intersection with c
-> no intersection with d

57 - 69
-> intersection with a -> 57-60 (-4) | 61-69 -> 53-56 | 61-69
-> no intersection with b
-> no intersection with c
-> no intersection with d

water_ranges = [Range { start: 81, len: 14 }, Range { start: 53, len: 4 }, Range { start: 61, len: 9 }]



water-to-light map:
a 88 18 7 -> 18 - 24 ; 88 - 94
b 18 25 70 -> 25 - 94 ; 18 - 87

81-94
-> no intersection a
-> intersection with b -> 81-94 -> 74-87

53-56
-> no a
-> yes b -> 53-56 -> 46-49

61-69
-> no a
-> yes b -> 61-69 -> 54-62

light_ranges = [Range { start: 74, len: 14 }, Range { start: 46, len: 4 }, Range { start: 54, len: 9 }]


light-to-temperature map:
a 45 77 23 -> 77-99 ; 45-67
b 81 45 19 -> 45-63 ; 81-99
c 68 64 13 -> 64-76 ; 68-80

74-87
-> yes a -> 74-76 | 77-87 (-32) -> 74-76 | 45-55
-> no b
-> yes c ->

46-49
54-62

temperature_ranges = [Range { start: 78, len: 3 }, Range { start: 81, len: 11 }, Range { start: 82, len: 4 }, Range { start: 90, len: 9 }]
