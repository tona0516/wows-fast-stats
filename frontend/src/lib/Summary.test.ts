import { Summary } from "src/lib/Summary";
import { domain } from "wailsjs/go/models";

test("undefined", () => {
  expect(
    Summary.calculate(undefined, [], new domain.UserConfig()),
  ).toBeUndefined();
});
