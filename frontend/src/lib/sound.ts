// lib/sound.ts

type SoundName =
  | "accepted"
  | "rejected"
  | "beep"
  | "winner"
  | "default"
  | "over";

class SoundManager {
  private enabled = false;
  private sounds: Record<SoundName, HTMLAudioElement>;
  private lastPlayedAt: Partial<Record<SoundName, number>> = {};
  private cooldownMs: Partial<Record<SoundName, number>> = {
    default: 200,
    accepted: 75,
    rejected: 75,
    winner: 250,
    over: 250,
    beep: 0,
  };

  constructor() {
    this.sounds = {
      accepted: new Audio("/sounds/accepted.wav"),
      rejected: new Audio("/sounds/rejected.wav"),
      beep: new Audio("/sounds/beep.wav"),
      winner: new Audio("/sounds/winner.wav"),
      default: new Audio("/sounds/default.wav"),
      over: new Audio("/sounds/finished.wav"),
    };

    Object.values(this.sounds).forEach((audio) => {
      audio.volume = 0.5;
    });
  }

  enable() {
    this.enabled = true;
  }

  disable() {
    this.enabled = false;
  }

  toggle() {
    this.enabled = !this.enabled;
    return this.enabled;
  }

  play(name: SoundName) {
    if (!this.enabled) {
      return;
    }

    const now = Date.now();
    const cooldown = this.cooldownMs[name] ?? 0;
    const last = this.lastPlayedAt[name] ?? 0;
    if (cooldown > 0 && now - last < cooldown) {
      return;
    }
    this.lastPlayedAt[name] = now;

    const sound = this.sounds[name];
    sound.currentTime = 0;
    sound.play().catch(() => {});
  }

  stop(name: SoundName) {
    if (!this.enabled) {
      return;
    }
    const sound = this.sounds[name];
    sound.pause();
    sound.currentTime = 0;
  }
}

export const soundManager = new SoundManager();
