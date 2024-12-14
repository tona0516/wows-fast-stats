import { storedExcludedPlayers } from "src/stores";

export namespace ExcludedPlayers {
  export const add = (accountID: number) =>
    storedExcludedPlayers.update((players) => {
      players.add(accountID);
      return players;
    });

  export const remove = (accountID: number) =>
    storedExcludedPlayers.update((players) => {
      players.delete(accountID);
      return players;
    });
}
