import { NumbersURL } from "src/lib/NumbersURL";
import { domain } from "wailsjs/go/models";

const isValid = (url: string): boolean => {
  try {
    new URL(url);
    return true;
  } catch (err) {
    return false;
  }
};

test("clan", () => {
  const player = new domain.Player({
    player_info: {
      clan: new domain.Clan({
        id: 1,
        tag: "TEST",
      }),
    },
  });

  expect(isValid(NumbersURL.clan(player))).toBe(true);
});

test("player", () => {
  const player = new domain.Player({
    player_info: {
      id: 1,
      name: "tonango",
    },
  });

  expect(isValid(NumbersURL.player(player))).toBe(true);
});

test("ship", () => {
  expect(isValid(NumbersURL.ship(1, "ARP Yamato"))).toBe(true);
});
