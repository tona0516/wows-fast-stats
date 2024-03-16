import { NumbersURL } from "src/lib/NumbersURL";

test("clan", () => {
  expect(NumbersURL.clan(1234567890)).toBe(
    "https://asia.wows-numbers.com/clan/1234567890,/",
  );
});

test("player", () => {
  expect(NumbersURL.player(1234567890, "test")).toBe(
    "https://asia.wows-numbers.com/player/1234567890,test/",
  );
});

test("ship", () => {
  expect(NumbersURL.ship(1234567890)).toBe(
    "https://asia.wows-numbers.com/ship/1234567890,/",
  );
});
