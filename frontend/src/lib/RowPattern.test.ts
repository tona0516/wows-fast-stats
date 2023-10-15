import { RowPattern } from "src/lib/RowPattern";
import { domain } from "wailsjs/go/models";

const STATS_PATTERN: string = "pvp_all";
const PLAYER = new domain.Player({
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
const ALL_COLUMN_NUMBER: number = 1;

test("no column", () => {
  expect(RowPattern.derive(PLAYER, STATS_PATTERN, 0)).toBe(
    RowPattern.NO_COLUMN,
  );
});

test("private", () => {
  const player = new domain.Player({
    player_info: {
      is_hidden: true,
    },
  });

  expect(RowPattern.derive(player, STATS_PATTERN, ALL_COLUMN_NUMBER)).toBe(
    RowPattern.PRIVATE,
  );
});

test("no data - 無効なアカウントID", () => {
  const player = new domain.Player({
    player_info: {
      is_hidden: false,
      id: 0,
    },
  });

  expect(RowPattern.derive(player, STATS_PATTERN, ALL_COLUMN_NUMBER)).toBe(
    RowPattern.NO_DATA,
  );
});

test("no data - 総合戦闘数=0", () => {
  const player = new domain.Player({
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

  expect(RowPattern.derive(player, STATS_PATTERN, ALL_COLUMN_NUMBER)).toBe(
    RowPattern.NO_DATA,
  );
});

test("no ship stats", () => {
  const player = new domain.Player({
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

  expect(RowPattern.derive(player, STATS_PATTERN, ALL_COLUMN_NUMBER)).toBe(
    RowPattern.NO_SHIP_STATS,
  );
});

test("full", () => {
  expect(RowPattern.derive(PLAYER, STATS_PATTERN, ALL_COLUMN_NUMBER)).toBe(
    RowPattern.FULL,
  );
});
