import { NumbersURL } from "src/lib/NumbersURL";

test("clan", () => {
  expect(NumbersURL.clan(1234567890)).toMatch(
    /^(https:\/\/asia\.wows-numbers\.com\/clan\/)([0-9]+,\/)$/,
  );
});

test("player", () => {
  expect(NumbersURL.player(1234567890)).toMatch(
    /^(https:\/\/asia\.wows-numbers\.com\/player\/)([0-9]+,\/)$/,
  );
});

test("ship", () => {
  expect(NumbersURL.ship(1234567890)).toMatch(
    /^(https:\/\/asia\.wows-numbers\.com\/ship\/)([0-9]+,\/)$/,
  );
});
