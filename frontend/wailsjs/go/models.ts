export namespace vo {
	
	export class Displays {
	    pr: boolean;
	    ship_damage: boolean;
	    ship_win_rate: boolean;
	    ship_kd_rate: boolean;
	    ship_win_survived_rate: boolean;
	    ship_lose_survived_rate: boolean;
	    ship_exp: boolean;
	    ship_battles: boolean;
	    player_damage: boolean;
	    player_win_rate: boolean;
	    player_kd_rate: boolean;
	    player_win_survived_rate: boolean;
	    player_lose_survived_rate: boolean;
	    player_exp: boolean;
	    player_battles: boolean;
	    player_avg_tier: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Displays(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pr = source["pr"];
	        this.ship_damage = source["ship_damage"];
	        this.ship_win_rate = source["ship_win_rate"];
	        this.ship_kd_rate = source["ship_kd_rate"];
	        this.ship_win_survived_rate = source["ship_win_survived_rate"];
	        this.ship_lose_survived_rate = source["ship_lose_survived_rate"];
	        this.ship_exp = source["ship_exp"];
	        this.ship_battles = source["ship_battles"];
	        this.player_damage = source["player_damage"];
	        this.player_win_rate = source["player_win_rate"];
	        this.player_kd_rate = source["player_kd_rate"];
	        this.player_win_survived_rate = source["player_win_survived_rate"];
	        this.player_lose_survived_rate = source["player_lose_survived_rate"];
	        this.player_exp = source["player_exp"];
	        this.player_battles = source["player_battles"];
	        this.player_avg_tier = source["player_avg_tier"];
	    }
	}
	export class PlayerPlayerStats {
	    battles: number;
	    avg_damage: number;
	    win_rate: number;
	    win_survived_rate: number;
	    lose_survived_rate: number;
	    kd_rate: number;
	    exp: number;
	    avg_tier: number;
	
	    static createFrom(source: any = {}) {
	        return new PlayerPlayerStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.avg_damage = source["avg_damage"];
	        this.win_rate = source["win_rate"];
	        this.win_survived_rate = source["win_survived_rate"];
	        this.lose_survived_rate = source["lose_survived_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.exp = source["exp"];
	        this.avg_tier = source["avg_tier"];
	    }
	}
	export class PlayerPlayerInfo {
	    id: number;
	    name: string;
	    clan: string;
	    is_hidden: boolean;
	    stats_url: string;
	
	    static createFrom(source: any = {}) {
	        return new PlayerPlayerInfo(source);
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
	export class PlayerShipStats {
	    battles: number;
	    avg_damage: number;
	    win_rate: number;
	    win_survived_rate: number;
	    lose_survived_rate: number;
	    kd_rate: number;
	    exp: number;
	    personal_rating: number;
	
	    static createFrom(source: any = {}) {
	        return new PlayerShipStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.battles = source["battles"];
	        this.avg_damage = source["avg_damage"];
	        this.win_rate = source["win_rate"];
	        this.win_survived_rate = source["win_survived_rate"];
	        this.lose_survived_rate = source["lose_survived_rate"];
	        this.kd_rate = source["kd_rate"];
	        this.exp = source["exp"];
	        this.personal_rating = source["personal_rating"];
	    }
	}
	export class PlayerShipInfo {
	    name: string;
	    nation: string;
	    tier: number;
	    type: string;
	    stats_url: string;
	
	    static createFrom(source: any = {}) {
	        return new PlayerShipInfo(source);
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
	    player_ship_info: PlayerShipInfo;
	    player_ship_stats: PlayerShipStats;
	    player_player_info: PlayerPlayerInfo;
	    player_player_stats: PlayerPlayerStats;
	
	    static createFrom(source: any = {}) {
	        return new Player(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.player_ship_info = this.convertValues(source["player_ship_info"], PlayerShipInfo);
	        this.player_ship_stats = this.convertValues(source["player_ship_stats"], PlayerShipStats);
	        this.player_player_info = this.convertValues(source["player_player_info"], PlayerPlayerInfo);
	        this.player_player_stats = this.convertValues(source["player_player_stats"], PlayerPlayerStats);
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
	
	
	
	
	export class TeamAverage {
	    personal_rating: number;
	    damage_by_ship: number;
	    win_rate_by_ship: number;
	    kd_rate_by_ship: number;
	    damage_by_player: number;
	    win_rate_by_player: number;
	    kd_rate_by_player: number;
	
	    static createFrom(source: any = {}) {
	        return new TeamAverage(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.personal_rating = source["personal_rating"];
	        this.damage_by_ship = source["damage_by_ship"];
	        this.win_rate_by_ship = source["win_rate_by_ship"];
	        this.kd_rate_by_ship = source["kd_rate_by_ship"];
	        this.damage_by_player = source["damage_by_player"];
	        this.win_rate_by_player = source["win_rate_by_player"];
	        this.kd_rate_by_player = source["kd_rate_by_player"];
	    }
	}
	export class Team {
	    players: Player[];
	    name: string;
	    team_average: TeamAverage;
	
	    static createFrom(source: any = {}) {
	        return new Team(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.players = this.convertValues(source["players"], Player);
	        this.name = source["name"];
	        this.team_average = this.convertValues(source["team_average"], TeamAverage);
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

