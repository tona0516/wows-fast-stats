import tonakoPointing from "src/assets/images/tonako_pointing.png";
import tonakoSorry from "src/assets/images/tonako_sorry.png";
import tonakoStandby from "src/assets/images/tonako_standby.png";

export class Tonako {
  private constructor(private readonly imgPath: string) {}

  static readonly Pointing = new Tonako(tonakoPointing);
  static readonly Standby = new Tonako(tonakoStandby);
  static readonly Sorry = new Tonako(tonakoSorry);

  getImgPath(): string {
    return this.imgPath;
  }
}
