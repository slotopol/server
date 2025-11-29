export interface Gamble {
  spin(mrtp: number): void;
  getBet(): number;
  setBet(bet: number): Promise<void>; // Assuming errors are handled via Promise rejections
}

export enum GT {
  SLOT = 1,
  KENO,
}

export enum GP {
  L_PAY = 1 << 0,
  R_PAY = 1 << 1,
  C_PAY = 1 << 2,
  S_PAY = 1 << 3,

  L_SEL = 1 << 4,
  W_SEL = 1 << 5,
  JACK = 1 << 6,
  FILL = 1 << 7,

  // Skip two
  CASC = 1 << 10,
  C_MULT = 1 << 11,

  FG_HAS = 1 << 12,
  RETRIG = 1 << 13,
  FG_REEL = 1 << 14,
  FG_MULT = 1 << 15,

  R_MULT = 1 << 16,
  SCAT = 1 << 17,
  WILD = 1 << 18,
  R_WILD = 1 << 19,

  B_WILD = 1 << 20,
  W_TURN = 1 << 21,
  W_MULT = 1 << 22,
  B_SYM = 1 << 23,

  FG_NO = 0,
}

export interface GameAlias {
  Prov: string;
  Name: string;
  Date?: number; // Unix timestamp
}

export interface AlgDescr {
  GT?: GT;
  GP?: GP;
  SX?: number;
  SY?: number;
  SN?: number;
  LN?: number;
  WN?: number;
  BN?: number;
  RTP: number[];
}

export interface AlgInfo {
  Aliases: GameAlias[];
  AlgDescr: AlgDescr;
}

export interface GameInfo {
  GameAlias: GameAlias;
  AlgDescr?: AlgDescr;
}

export type Scanner = (context: any, mrtp: number) => number; // Assuming context is 'any' for now

export const GameFactory: { [key: string]: () => Gamble } = {};
export const ScanFactory: { [key: string]: Scanner } = {};
export const InfoMap: { [key: string]: GameInfo } = {};
export const AlgList: AlgInfo[] = [];

// Helper function (Date and Year are assumed to be available if needed)
export const MakeRtpList = <T>(reelsMap: { [key: number]: T }): number[] => {
  const list = Object.keys(reelsMap).map(Number);
  list.sort((a, b) => a - b);
  return list;
};

export const Date = (year: number, month: number, day: number): number => {
    return new Date(year, month -1, day).getTime() / 1000;
}

export class AlgInfoImpl implements AlgInfo {
  Aliases: GameAlias[];
  AlgDescr: AlgDescr;

  constructor(aliases: GameAlias[], algDescr: AlgDescr) {
    this.Aliases = aliases;
    this.AlgDescr = algDescr;
  }

  SetupFactory(game: () => Gamble, scan?: Scanner) {
    AlgList.push(this);
    for (const ga of this.Aliases) {
      const aid = `${ga.Prov}/${ga.Name}`.toLowerCase().replace(/[^a-z0-9/]/g, '');
      InfoMap[aid] = {
        GameAlias: ga,
        AlgDescr: this.AlgDescr,
      };
      GameFactory[aid] = game;
      if (scan) {
        ScanFactory[aid] = scan;
      }
    }
  }
}
