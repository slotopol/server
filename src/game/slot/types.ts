export type Sym = number; // byte
export type Pos = number; // int8

/**
 * Represents the y-coordinates of a payline across the reels.
 * The index of the array corresponds to the reel number (0-indexed),
 * and the value at that index is the y-position of the symbol on that reel.
 * For example, `[0, 1, 2, 1, 0]` means the line passes through:
 * reel 0, y-pos 0
 * reel 1, y-pos 1
 * reel 2, y-pos 2
 * reel 3, y-pos 1
 * reel 4, y-pos 0
 */
export type Linex = Pos[];

export interface Reels {
  cols(): number;
  reel(col: Pos): Sym[];
  reshuffles(): number; // uint64
  toString(): string;
}

// Reels5x is an array of 5 reels, where each reel is an array of symbols.
export type Reels5x = Sym[][];

export interface WinItem {
  pay?: number;
  mult?: number;
  sym?: Sym;
  num?: Pos;
  line?: number;
  xy: Linex; // symbols positions on screen
  free?: number;
  bid?: number; // bonus identifier
  bon?: any; // bonus game data - Go type is `any`
  jid?: number; // jackpot identifier
  jack?: number; // jackpot win
}

export type Wins = WinItem[];

export interface SlotGame {
  clone(): SlotGame;
  scanner(wins: Wins): Error | null; // Appends to wins
  cost(): { cost: number; hasJackpotRate: boolean };
  free(): boolean;
  spin(rtp: number): void;
  spawn(wins: Wins, fund: number, rtp: number): void;
  prepare(): void;
  apply(wins: Wins): void;
  getGain(): number;
  setGain(gain: number): Error | null;
  getBet(): number;
  setBet(bet: number): Error | null;
  getSel(): number;
  setSel(sel: number): Error | null;
  setMode(mode: number): Error | null;
}

export interface Slotx {
  sel: number; // selected bet lines
  bet: number; // bet value
  gain?: number; // gain for double up games
  fsn?: number; // free spin number
  fsr?: number; // free spin remains

  // Methods from the Go Slotx struct
  cost(): { cost: number; hasJackpotRate: boolean };
  free(): boolean;
  spawn(wins: Wins, fund: number, rtp: number): void;
  prepare(): void;
  apply(wins: Wins): void;
  getGain(): number;
  setGain(gain: number): Error | null;
  getBet(): number;
  setBet(bet: number): Error | null;
  getSel(): number;
  setSelNum(sel: number, bln: number): Error | null; // Renamed from SetSel to avoid conflict with SlotGame
  setMode(mode: number): Error | null;
}

export interface Screen {
  dim(): { sx: Pos; sy: Pos };
  at(x: Pos, y: Pos): Sym;
  ly(x: Pos, line: Linex): Sym;
  setSym(x: Pos, y: Pos, sym: Sym): void;
  setCol(x: Pos, reel: Sym[], pos: number): void;
  reelSpin(reels: Reels): void;
  symNum(sym: Sym): Pos;
  scatNum(scat: Sym): Pos;
  scatPos(scat: Sym): Linex;
}

export class Screen5x3 implements Screen {
  public scr: Sym[][]; // 5 columns, 3 rows

  constructor() {
    this.scr = Array(5)
      .fill(null)
      .map(() => Array(3).fill(0));
  }

  dim(): { sx: Pos; sy: Pos } {
    return { sx: 5, sy: 3 };
  }

  at(x: Pos, y: Pos): Sym {
    // Adjust for 0-indexed arrays from 1-indexed Go version
    return this.scr[x - 1][y - 1];
  }

  ly(x: Pos, line: Linex): Sym {
    // Adjust for 0-indexed arrays
    return this.scr[x - 1][line[x - 1] - 1];
  }

  setSym(x: Pos, y: Pos, sym: Sym): void {
    // Adjust for 0-indexed arrays
    this.scr[x - 1][y - 1] = sym;
  }

  setCol(x: Pos, reel: Sym[], pos: number): void {
    // Adjust for 0-indexed arrays
    const col = x - 1;
    for (let y = 0; y < 3; y++) {
      this.scr[col][y] = reel[(pos + y) % reel.length];
    }
  }

  reelSpin(reels: Reels): void {
    for (let x = 1; x <= 5; x++) {
      const reel = reels.reel(x as Pos);
      const hit = Math.floor(Math.random() * reel.length);
      this.setCol(x as Pos, reel, hit);
    }
  }

  symNum(sym: Sym): Pos {
    let n: Pos = 0;
    for (let x = 0; x < 5; x++) {
      for (let y = 0; y < 3; y++) {
        if (this.scr[x][y] === sym) {
          n++;
        }
      }
    }
    return n;
  }

  scatNum(scat: Sym): Pos {
    let n: Pos = 0;
    for (let x = 0; x < 5; x++) {
      if (this.scr[x][0] === scat || this.scr[x][1] === scat || this.scr[x][2] === scat) {
        n++;
      }
    }
    return n;
  }

  scatPos(scat: Sym): Linex {
    const line: Linex = Array(5).fill(0);
    for (let x = 0; x < 5; x++) {
      if (this.scr[x][0] === scat) {
        line[x] = 1;
      } else if (this.scr[x][1] === scat) {
        line[x] = 2;
      } else if (this.scr[x][2] === scat) {
        line[x] = 3;
      }
    }
    return line;
  }
}

export class Stat {
  private planned: number = 0; // uint64
  private reshuffles: number[] = Array(10).fill(0); // uint64 array
  private errcount: number = 0; // uint64
  private linepay: number = 0;
  private scatpay: number = 0;
  private freecount: number = 0; // uint64
  private freehits: number = 0; // uint64
  private bonuscount: number[] = Array(8).fill(0); // uint64 array
  private jackcount: number[] = Array(4).fill(0); // uint64 array

  // Mutexes are not directly translatable to typical browser/Node.js TS.
  // If concurrency control is needed, it would be handled differently.
  // For now, direct access is used.

  setPlan(n: number): void {
    this.planned = n;
  }

  getPlanned(): number {
    return this.planned;
  }

  getCount(): number {
    return this.reshuffles[0] - this.errcount;
  }

  getReshuf(cfn: number): number {
    return this.reshuffles[cfn - 1];
  }

  incErr(): void {
    this.errcount++;
  }

  lineRTP(cost: number): number {
    const reshuf = this.reshuffles[0] - this.errcount;
    if (reshuf === 0 || cost === 0) return 0;
    return (this.linepay / reshuf / cost) * 100;
  }

  scatRTP(cost: number): number {
    const reshuf = this.reshuffles[0] - this.errcount;
    if (reshuf === 0 || cost === 0) return 0;
    return (this.scatpay / reshuf / cost) * 100;
  }

  getFreeCount(): number {
    return this.freecount;
  }

  getFreeHits(): number {
    return this.freehits;
  }

  getBonusCount(bid: number): number {
    // Assuming bid is 1-indexed from Go, adjust if necessary
    return this.bonuscount[bid > 0 ? bid -1 : 0];
  }

  getJackCount(jid: number): number {
    // Assuming jid is 1-indexed from Go, adjust if necessary
    return this.jackcount[jid > 0 ? jid -1 : 0];
  }

  update(wins: Wins, cfn: number): void {
    for (const wi of wins) {
      if (wi.pay && wi.pay !== 0) {
        const payment = wi.pay * (wi.mult || 1);
        if (wi.line && wi.line !== 0) {
          this.linepay += payment;
        } else {
          this.scatpay += payment;
        }
      }
      if (wi.free && wi.free !== 0) {
        this.freecount += wi.free;
        this.freehits++;
      }
      if (wi.bid && wi.bid !== 0 && wi.bid > 0 && wi.bid <= this.bonuscount.length) {
        this.bonuscount[wi.bid -1]++;
      }
      if (wi.jid && wi.jid !== 0 && wi.jid > 0 && wi.jid <= this.jackcount.length) {
        this.jackcount[wi.jid-1]++;
      }
    }
    if (cfn > 0 && cfn <= this.reshuffles.length) {
      this.reshuffles[cfn - 1]++;
    }
  }
}

// Errors (can be simple Error objects or custom error classes)
export const ErrNoWay = new Error("no way to here");
export const ErrBetEmpty = new Error("bet is empty");
export const ErrNoLineset = new Error("lines set is empty");
export const ErrLinesetOut = new Error("lines set is out of range bet lines");
export const ErrNoFeature = new Error("feature not available");
export const ErrDisabled = new Error("feature is disabled");

export const JackBasis = 250_000_000;

// Helper function to mimic slot.FindClosest (simplified: finds first match or default)
export function findClosestReels(reelsMap: Record<number, Reels5x>, targetRtp: number): Reels5x | undefined {
  // Convert string keys to numbers for comparison if necessary, though YAML keys are numbers
  const rtps = Object.keys(reelsMap).map(Number).sort((a, b) => Math.abs(a - targetRtp) - Math.abs(b - targetRtp));
  if (rtps.length > 0) {
    return reelsMap[rtps[0]];
  }
  // Fallback: return the first available reel set if no specific match logic is implemented
  const firstKey = Object.keys(reelsMap)[0];
  return firstKey ? reelsMap[Number(firstKey)] : undefined;
}
// Helper function to mimic slot.ReadObj (simplified: returns the object directly)
export function readObj<T>(data: T): T {
    return data;
}

// Helper function to mimic slot.ReadMap (simplified: returns the map directly)
export function readMap<T>(data: Record<string, T>): Record<string, T> {
    return data;
}
