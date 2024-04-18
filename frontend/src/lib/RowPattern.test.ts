import { RowPattern } from "src/lib/RowPattern";
import { data } from "wailsjs/go/models";

const STATS_PATTERN: string = "pvp_all";
const PLAYER = new data.Player({
  player_info: {
    is_hidden: false,
    id: 1,
  },
  pvp_all: {
    ship: {
      battles: 1,
    },
    overall: {
      battles: 1,
    },
  },
});
const ALL_COLUMN_COUNT: number = 1;
const SHIP_COLUMN_COUNT: number = 1;

test("no column", () => {
  expect(RowPattern.derive(PLAYER, STATS_PATTERN, 0, 0)).toBe(
    RowPattern.NO_COLUMN,
  );
});

test("private", () => {
  const player = new data.Player({
    player_info: {
      is_hidden: true,
    },
  });

  expect(
    RowPattern.derive(
      player,
      STATS_PATTERN,
      ALL_COLUMN_COUNT,
      SHIP_COLUMN_COUNT,
    ),
  ).toBe(RowPattern.PRIVATE);
});

test("no data - 無効なアカウントID", () => {
  const player = new data.Player({
    player_info: {
      is_hidden: false,
      id: 0,
    },
  });

  expect(
    RowPattern.derive(
      player,
      STATS_PATTERN,
      ALL_COLUMN_COUNT,
      SHIP_COLUMN_COUNT,
    ),
  ).toBe(RowPattern.NO_STATS);
});

test("no data - 総合戦闘数=0", () => {
  const player = new data.Player({
    player_info: {
      is_hidden: false,
      id: 1,
    },
    pvp_all: {
      overall: {
        battles: 0,
      },
    },
  });

  expect(
    RowPattern.derive(
      player,
      STATS_PATTERN,
      ALL_COLUMN_COUNT,
      SHIP_COLUMN_COUNT,
    ),
  ).toBe(RowPattern.NO_STATS);
});

test("no ship stats", () => {
  const player = new data.Player({
    player_info: {
      is_hidden: false,
      id: 1,
    },
    pvp_all: {
      ship: {
        battles: 0,
      },
      overall: {
        battles: 1,
      },
    },
  });

  expect(
    RowPattern.derive(
      player,
      STATS_PATTERN,
      ALL_COLUMN_COUNT,
      SHIP_COLUMN_COUNT,
    ),
  ).toBe(RowPattern.NO_SHIP_STATS);
});

test("full", () => {
  expect(
    RowPattern.derive(
      PLAYER,
      STATS_PATTERN,
      ALL_COLUMN_COUNT,
      SHIP_COLUMN_COUNT,
    ),
  ).toBe(RowPattern.FULL);
  expect(RowPattern.derive(PLAYER, STATS_PATTERN, ALL_COLUMN_COUNT, 0)).toBe(
    RowPattern.FULL,
  );
});
