import { SlotGame, Reels5x } from './slot/types'; // Assuming SlotGame is defined here

/**
 * GameType typically represents the category of the game.
 * GTslot would be one of these.
 */
export enum GameType {
  GTslot = 1,
  // Add other game types if known, e.g., GTkeno, GTroulette
}

/**
 * GameProps are bitmask flags representing properties/features of the game.
 * Example values from flowers_link.go:
 * GPlpay | GPlsel | GPretrig | GPfgreel | GPfgmult | GPwild
 * These would be individual enum members.
 */
export enum GameProps {
  GPlpay = 1 << 0, // Has line pay
  GPlsel = 1 << 1, // Has selectable lines
  GPretrig = 1 << 2, // Has retriggerable free games
  GPfgreel = 1 << 3, // Has different reels for free games
  GPfgmult = 1 << 4, // Has multiplier for free games
  GPwild = 1 << 5, // Has wild symbols
  // Add other properties as they are discovered or needed
  // For example, from the `gp` value in the GitHub README: 4628497
  // This implies more flags exist. For now, only include those from flowers.
  GPscatter = 1 << 6, // Has scatter symbols (example, if needed)
  GPbonus = 1 << 7,   // Has a bonus game (example, if needed)
}

/**
 * GameAlias represents an alternative name or version for a game.
 */
export interface GameAlias {
  prov: string; // Provider (e.g., "NetEnt")
  name: string; // Game Name (e.g., "Flowers")
  date?: GameDate; // Release date (from game.Date)
  year?: number; // Fallback if only year is known
}

/**
 * AlgDescr describes the algorithm features of a game.
 */
export interface AlgDescr {
  gt: GameType; // Game Type (e.g., GTslot)
  gp: number; // Game Properties (bitmask of GameProps)
  sx: number; // Screen X dimension (e.g., 5 reels)
  sy: number; // Screen Y dimension (e.g., 3 rows)
  sn: number; // Number of symbols in LinePay (max symbol ID?)
  ln: number; // Number of bet lines
  bn?: number; // Number of bonus games or types (0 if none)
  rtp: RTPInfo[]; // List of available RTP configurations
}

/**
 * RTPInfo describes a specific RTP configuration for a game.
 * In the Go code, it's a list like: [87.788791, 89.230191, ...]
 * It can also be a map like in flowers_reel.yaml where key is RTP and value is Reels5x
 * For flexibility, we can define it to support a list of numbers or a map.
 * Given `game.MakeRtpList(ReelsMap)` in flowers_link.go, it seems to be a list of numbers.
 */
export interface RTPInfo {
  value: number; // The RTP percentage value
  // reels?: Reels5x; // Optional: if specific reels are tied to this RTP entry directly
}


/**
 * AlgInfo is a comprehensive structure holding all information about a game algorithm.
 */
export interface AlgInfo {
  aliases: GameAlias[];
  algDescr: AlgDescr;
  // The Go code has a SetupFactory method.
  // In TypeScript, this could be a factory function or a class constructor.
  gameFactory?: () => Gamble; // Factory function to create a new game instance
  calcStat?: (rtp: number) => number; // Reference to the main stat calculation function
}

/**
 * Gamble is an alias for SlotGame, representing the game instance.
 */
export type Gamble = SlotGame;

/**
 * GameDate represents a date. The Go `game.Date(YYYY, MM, DD)` suggests it could be simple.
 * Using a string in ISO format or a simple object.
 */
export type GameDate = string; // e.g., "2013-11-11" or  number (timestamp) or  { year: number; month: number; day: number };

/**
 * Helper function to mimic game.MakeRtpList(ReelsMap)
 * It extracts the RTP values (keys) from a ReelsMap and returns them as a list of RTPInfo objects.
 */
export function makeRtpList(reelsMap: Record<number, Reels5x>): RTPInfo[] {
  return Object.keys(reelsMap).map(rtpStr => ({ value: parseFloat(rtpStr) })).sort((a,b) => a.value - b.value);
}

/**
 * Represents the structure for game registration.
 */
export interface GameRegistration {
  info: AlgInfo;
  factory: () => Gamble;
  statCalc: (ctx: any, mrtp: number) => number; // Using 'any' for context for now
}

const registeredGames: Map<string, GameRegistration> = new Map();

export function registerGame(gameName: string, registration: GameRegistration): void {
  if (registeredGames.has(gameName)) {
    console.warn(`Game ${gameName} is already registered. Overwriting.`);
  }
  registeredGames.set(gameName, registration);
}

export function getGameRegistration(gameName: string): GameRegistration | undefined {
  return registeredGames.get(gameName);
}
