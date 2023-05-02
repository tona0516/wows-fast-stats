/**
 * private: hidden player.
 * nodata: invalid player(bot/deleted account) or 0 for all random battle.
 * noshipstats: 0 in for random battle with the ship.
 * nopr: not exist expected value in numbers api.
 * full: all values exists.
 */
type DisplayPattern = "private" | "nodata" | "noshipstats" | "nopr" | "full";
