export namespace vo {
	
	export class Basic {
	    player_name: boolean;
	    ship_info: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Basic(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.player_name = source["player_name"];
	        this.ship_info = source["ship_info"];
	    }
	}
	export class TierGroup[float64] {
	    low: number;
	    middle: number;
	    high: number;
	
	    static createFrom(source: any = {}) {
	        return new TierGroup[float64](source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.low = source["low"];
	        this.middle = source["middle"];
	        this.high = source["high"];
	    }
	}
	export class ShipTypeGroup[float64] {
	    ss: number;
	    dd: number;
	    cl: number;
	    bb: number;
	    cv: number;
	
	    static createFrom(source: any = {}) {
	        return new ShipTypeGroup[float64](source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ss = source["ss"];
	        this.dd = source["dd"];
	        this.cl = source["cl"];
	        this.bb = source["bb"];
	        this.cv = source["cv"];
	    }
	}
	export class PlayerStats {
	    battles: number;
	    damage: number;
	    win_rate: number;
	    win_survived_rate: number;
	    lose_survived_rate: number;
	    kd_rate: number;
	    exp: number;
	    avg_tier: number;
	    using_ship_type_rate: ShipTypeGroup[float64];
	    using_tier_rate: TierGroup[float64];
	
	    static createFrom(source: any = {}) {
	        return new PlayerStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.damage = source["damage"];
	        this.win_rate = source["win_rate"];
	        this.win_survived_rate = source["win_survived_rate"];
	        this.lose_survived_rate = source["lose_survived_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.exp = source["exp"];
	        this.avg_tier = source["avg_tier"];
	        this.using_ship_type_rate = this.convertValues(source["using_ship_type_rate"], ShipTypeGroup[float64]);
	        this.using_tier_rate = this.convertValues(source["using_tier_rate"], TierGroup[float64]);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PlayerInfo {
	    id: number;
	    name: string;
	    clan: string;
	    is_hidden: boolean;
	    stats_url: string;
	
	    static createFrom(source: any = {}) {
	        return new PlayerInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.clan = source["clan"];
	        this.is_hidden = source["is_hidden"];
	        this.stats_url = source["stats_url"];
	    }
	}
	export class ShipStats {
	    battles: number;
	    damage: number;
	    win_rate: number;
	    win_survived_rate: number;
	    lose_survived_rate: number;
	    kd_rate: number;
	    exp: number;
	    pr: number;
	
	    static createFrom(source: any = {}) {
	        return new ShipStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.damage = source["damage"];
	        this.win_rate = source["win_rate"];
	        this.win_survived_rate = source["win_survived_rate"];
	        this.lose_survived_rate = source["lose_survived_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.exp = source["exp"];
	        this.pr = source["pr"];
	    }
	}
	export class ShipInfo {
	    name: string;
	    nation: string;
	    tier: number;
	    type: string;
	    stats_url: string;
	
	    static createFrom(source: any = {}) {
	        return new ShipInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.nation = source["nation"];
	        this.tier = source["tier"];
	        this.type = source["type"];
	        this.stats_url = source["stats_url"];
	    }
	}
	export class Player {
	    ship_info: ShipInfo;
	    ship_stats: ShipStats;
	    player_info: PlayerInfo;
	    player_stats: PlayerStats;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ship_info = this.convertValues(source["ship_info"], ShipInfo);
	        this.ship_stats = this.convertValues(source["ship_stats"], ShipStats);
	        this.player_info = this.convertValues(source["player_info"], PlayerInfo);
	        this.player_stats = this.convertValues(source["player_stats"], PlayerStats);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Team {
	    players: Player[];
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Team(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.players = this.convertValues(source["players"], Player);
	        this.name = source["name"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class OverallComp {
	    damage: Between;
	    win_rate: Between;
	    kd_rate: Between;
	
	    static createFrom(source: any = {}) {
	        return new OverallComp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.damage = this.convertValues(source["damage"], Between);
	        this.win_rate = this.convertValues(source["win_rate"], Between);
	        this.kd_rate = this.convertValues(source["kd_rate"], Between);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Between {
	    friend: number;
	    enemy: number;
	    diff: number;
	
	    static createFrom(source: any = {}) {
	        return new Between(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.friend = source["friend"];
	        this.enemy = source["enemy"];
	        this.diff = source["diff"];
	    }
	}
	export class ShipComp {
	    pr: Between;
	    damage: Between;
	    win_rate: Between;
	    kd_rate: Between;
	
	    static createFrom(source: any = {}) {
	        return new ShipComp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pr = this.convertValues(source["pr"], Between);
	        this.damage = this.convertValues(source["damage"], Between);
	        this.win_rate = this.convertValues(source["win_rate"], Between);
	        this.kd_rate = this.convertValues(source["kd_rate"], Between);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Comparision {
	    ship: ShipComp;
	    overall: OverallComp;
	
	    static createFrom(source: any = {}) {
	        return new Comparision(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ship = this.convertValues(source["ship"], ShipComp);
	        this.overall = this.convertValues(source["overall"], OverallComp);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Meta {
	    date: string;
	    arena: string;
	    type: string;
	    own_ship: string;
	
	    static createFrom(source: any = {}) {
	        return new Meta(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.date = source["date"];
	        this.arena = source["arena"];
	        this.type = source["type"];
	        this.own_ship = source["own_ship"];
	    }
	}
	export class Battle {
	    meta: Meta;
	    comparision: Comparision;
	    teams: Team[];
	
	    static createFrom(source: any = {}) {
	        return new Battle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.meta = this.convertValues(source["meta"], Meta);
	        this.comparision = this.convertValues(source["comparision"], Comparision);
	        this.teams = this.convertValues(source["teams"], Team);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class Overall {
	    damage: boolean;
	    win_rate: boolean;
	    kd_rate: boolean;
	    win_survived_rate: boolean;
	    lose_survived_rate: boolean;
	    exp: boolean;
	    battles: boolean;
	    avg_tier: boolean;
	    using_ship_type_rate: boolean;
	    using_tier_rate: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Overall(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.damage = source["damage"];
	        this.win_rate = source["win_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.win_survived_rate = source["win_survived_rate"];
	        this.lose_survived_rate = source["lose_survived_rate"];
	        this.exp = source["exp"];
	        this.battles = source["battles"];
	        this.avg_tier = source["avg_tier"];
	        this.using_ship_type_rate = source["using_ship_type_rate"];
	        this.using_tier_rate = source["using_tier_rate"];
	    }
	}
	export class Ship {
	    pr: boolean;
	    damage: boolean;
	    win_rate: boolean;
	    kd_rate: boolean;
	    win_survived_rate: boolean;
	    lose_survived_rate: boolean;
	    exp: boolean;
	    battles: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Ship(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pr = source["pr"];
	        this.damage = source["damage"];
	        this.win_rate = source["win_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.win_survived_rate = source["win_survived_rate"];
	        this.lose_survived_rate = source["lose_survived_rate"];
	        this.exp = source["exp"];
	        this.battles = source["battles"];
	    }
	}
	export class Displays {
	    basic: Basic;
	    ship: Ship;
	    overall: Overall;
	
	    static createFrom(source: any = {}) {
	        return new Displays(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.basic = this.convertValues(source["basic"], Basic);
	        this.ship = this.convertValues(source["ship"], Ship);
	        this.overall = this.convertValues(source["overall"], Overall);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	
	
	
	
	
	
	
	
	
	
	
	export class UserConfig {
	    install_path: string;
	    appid: string;
	    font_size: string;
	    displays: Displays;
	    save_screenshot: boolean;
	    save_temp_arena_info: boolean;
	
	    static createFrom(source: any = {}) {
	        return new UserConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.install_path = source["install_path"];
	        this.appid = source["appid"];
	        this.font_size = source["font_size"];
	        this.displays = this.convertValues(source["displays"], Displays);
	        this.save_screenshot = source["save_screenshot"];
	        this.save_temp_arena_info = source["save_temp_arena_info"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

