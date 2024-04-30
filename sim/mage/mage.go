package mage

import (
	"time"

	"github.com/wowsims/cata/sim/core"
	"github.com/wowsims/cata/sim/core/proto"
	"github.com/wowsims/cata/sim/core/stats"
)

const (
	SpellFlagMage     = core.SpellFlagAgentReserved1
	Unused2           = core.SpellFlagAgentReserved2
	HotStreakSpells   = core.SpellFlagAgentReserved3
	BrainFreezeSpells = core.SpellFlagAgentReserved4
)

var TalentTreeSizes = [3]int{21, 21, 19}

type Mage struct {
	core.Character

	PresenceOfMindMod  *core.SpellMod
	arcanePowerCostMod *core.SpellMod
	arcanePowerDmgMod  *core.SpellMod
	arcanePowerGCDmod  *core.SpellMod

	Talents       *proto.MageTalents
	Options       *proto.MageOptions
	ArcaneOptions *proto.ArcaneMage_Options
	FireOptions   *proto.FireMage_Options
	FrostOptions  *proto.FrostMage_Options

	mirrorImage  *MirrorImage
	flameOrb     *FlameOrb
	frostfireOrb *FrostfireOrb

	ArcaneBarrage           *core.Spell
	ArcaneBlast             *core.Spell
	ArcaneExplosion         *core.Spell
	ArcaneMissiles          *core.Spell
	ArcaneMissilesTickSpell *core.Spell
	ArcanePower             *core.Spell
	Blizzard                *core.Spell
	Combustion              *core.Spell
	DeepFreeze              *core.Spell
	Ignite                  *core.Spell
	IgniteImpact            *core.Spell
	LivingBomb              *core.Spell
	LivingBombImpact        *core.Spell
	Fireball                *core.Spell
	FireBlast               *core.Spell
	FlameOrb                *core.Spell
	FlameOrbExplode         *core.Spell
	Flamestrike             *core.Spell
	Freeze                  *core.Spell
	Frostbolt               *core.Spell
	FrostfireBolt           *core.Spell
	FrostfireOrb            *core.Spell
	FrostfireOrbTickSpell   *core.Spell
	IceLance                *core.Spell
	PresenceOfMind          *core.Spell
	Pyroblast               *core.Spell
	PyroblastDot            *core.Spell
	PyroblastDotImpact      *core.Spell
	SummonWaterElemental    *core.Spell
	Scorch                  *core.Spell
	MirrorImage             *core.Spell
	BlastWave               *core.Spell
	DragonsBreath           *core.Spell
	IcyVeins                *core.Spell

	ArcaneBlastAura        *core.Aura
	ArcaneMissilesProcAura *core.Aura
	ArcanePotencyAura      *core.Aura
	ArcanePowerAura        *core.Aura
	BrainFreezeAura        *core.Aura
	ClearcastingAura       *core.Aura
	CriticalMassAuras      core.AuraArray
	FingersOfFrostAura     *core.Aura
	FrostArmorAura         *core.Aura
	GlyphedFrostArmorPA    *core.PendingAction
	hotStreakCritAura      *core.Aura
	HotStreakAura          *core.Aura
	IgniteDamageTracker    *core.Aura
	ImpactAura             *core.Aura
	PresenceOfMindAura     *core.Aura
	PyromaniacAura         *core.Aura

	ScalingBaseDamage float64

	CritDebuffCategories core.ExclusiveCategoryArray
}

func (mage *Mage) GetCharacter() *core.Character {
	return &mage.Character
}

func (mage *Mage) GetMage() *Mage {
	return mage
}

func (mage *Mage) HasPrimeGlyph(glyph proto.MagePrimeGlyph) bool {
	return mage.HasGlyph(int32(glyph))
}

func (mage *Mage) HasMajorGlyph(glyph proto.MageMajorGlyph) bool {
	return mage.HasGlyph(int32(glyph))
}
func (mage *Mage) HasMinorGlyph(glyph proto.MageMinorGlyph) bool {
	return mage.HasGlyph(int32(glyph))
}

func (mage *Mage) AddRaidBuffs(raidBuffs *proto.RaidBuffs) {
	raidBuffs.ArcaneBrilliance = true

	// if mage.Talents.ArcaneEmpowerment == 3 {
	// 	raidBuffs.ArcaneEmpowerment = true
	// }
}
func (mage *Mage) AddPartyBuffs(partyBuffs *proto.PartyBuffs) {
}

func (mage *Mage) ApplyTalents() {
	mage.ApplyArmorSpecializationEffect(stats.Intellect, proto.ArmorType_ArmorTypeCloth)

	mage.ApplyArcaneTalents()
	mage.ApplyFireTalents()
	mage.ApplyFrostTalents()

	mage.applyGlyphs()
}

func (mage *Mage) Initialize() {

	mage.applyArmorSpells()
	mage.registerArcaneBlastSpell()
	mage.registerArcaneExplosionSpell()
	mage.registerArcaneMissilesSpell()
	mage.registerBlizzardSpell()
	mage.registerDeepFreezeSpell()
	mage.registerFireballSpell()
	mage.registerFireBlastSpell()
	mage.registerFlameOrbSpell()
	mage.registerFlameOrbExplodeSpell()
	mage.registerFlamestrikeSpell()
	mage.registerFreezeSpell()
	mage.registerFrostboltSpell()
	mage.registerFrostfireOrbSpell()
	mage.registerIceLanceSpell()
	mage.registerScorchSpell()
	mage.registerLivingBombSpell()
	mage.registerFrostfireBoltSpell()
	mage.registerEvocation()
	mage.registerManaGemsCD()
	mage.registerMirrorImageCD()
	mage.registerCombustionSpell()
	mage.registerBlastWaveSpell()
	mage.registerDragonsBreathSpell()
	// mage.registerSummonWaterElementalCD()

	mage.applyArcaneMissileProc()
	mage.ScalingBaseDamage = 937.330078125
}

func (mage *Mage) Reset(sim *core.Simulation) {
}

func NewMage(character *core.Character, options *proto.Player, mageOptions *proto.MageOptions) *Mage {
	mage := &Mage{
		Character: *character,
		Talents:   &proto.MageTalents{},
		Options:   mageOptions,
	}

	core.FillTalentsProto(mage.Talents.ProtoReflect(), options.TalentsString, TalentTreeSizes)

	mage.mirrorImage = mage.NewMirrorImage()
	mage.flameOrb = mage.NewFlameOrb()
	mage.frostfireOrb = mage.NewFrostfireOrb()
	mage.EnableManaBar()
	return mage
}

func (mage *Mage) applyArmorSpells() {
	// Molten Armor
	// +3% spell crit, +5% with glyph

	critToAdd := 3 * core.CritRatingPerCritChance
	if mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfMoltenArmor) {
		critToAdd = 5 * core.CritRatingPerCritChance
	}

	mageArmorEffectCategory := "MageArmors"

	moltenArmor := mage.RegisterAura(core.Aura{
		Label:    "Molten Armor",
		ActionID: core.ActionID{SpellID: 30482},
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCrit, critToAdd)
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			mage.AddStatDynamic(sim, stats.SpellCrit, -critToAdd)
		},
	})

	moltenArmor.NewExclusiveEffect(mageArmorEffectCategory, true, core.ExclusiveEffect{})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 30482},
		SpellSchool:    core.SpellSchoolFire,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellMoltenArmor,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !moltenArmor.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			moltenArmor.Activate(sim)
		},
	})

	// Mage Armor
	// Restores 3% of your max mana every 5 seconds (+20% affect with glyph)
	mageArmorManaMetric := mage.NewManaMetrics(core.ActionID{SpellID: 6117})
	hasGlyph := mage.HasPrimeGlyph(proto.MagePrimeGlyph_GlyphOfMageArmor)
	manaRegenPer5Second := core.TernaryFloat64(hasGlyph, .036, 0.03)

	var pa *core.PendingAction
	mageArmor := mage.RegisterAura(core.Aura{
		ActionID: core.ActionID{SpellID: 6117},
		Label:    "Mage Armor",
		Duration: core.NeverExpires,
		OnGain: func(aura *core.Aura, sim *core.Simulation) {
			pa = core.StartPeriodicAction(sim, core.PeriodicActionOptions{
				Period: time.Second * 5,
				OnAction: func(sim *core.Simulation) {
					mage.AddMana(sim, mage.MaxMana()*manaRegenPer5Second, mageArmorManaMetric)
				},
			})
		},
		OnExpire: func(aura *core.Aura, sim *core.Simulation) {
			pa.Cancel(sim)
		},
	})

	mageArmor.NewExclusiveEffect(mageArmorEffectCategory, true, core.ExclusiveEffect{})

	mage.RegisterSpell(core.SpellConfig{
		ActionID:       core.ActionID{SpellID: 6117},
		SpellSchool:    core.SpellSchoolArcane,
		Flags:          core.SpellFlagAPL,
		ClassSpellMask: MageSpellMageArmor,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				GCD: core.GCDDefault,
			},
		},
		ExtraCastCondition: func(sim *core.Simulation, target *core.Unit) bool {
			return !mageArmor.IsActive()
		},

		ApplyEffects: func(sim *core.Simulation, _ *core.Unit, _ *core.Spell) {
			mageArmor.Activate(sim)
		},
	})
}

// Agent is a generic way to access underlying mage on any of the agents.
type MageAgent interface {
	GetMage() *Mage
}

func (mage *Mage) applyArcaneMissileProc() {
	if mage.Talents.HotStreak || mage.Talents.BrainFreeze > 0 {
		return
	}

	// Aura for when proc is successful
	mage.ArcaneMissilesProcAura = mage.RegisterAura(core.Aura{
		Label:    "Arcane Missiles Proc",
		ActionID: core.ActionID{SpellID: 79683},
		Duration: time.Second * 20,
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if spell == mage.ArcaneMissiles {
				aura.Deactivate(sim)
			}
		},
	})

	procChance := 0.4

	// Listener for procs
	mage.RegisterAura(core.Aura{
		Label:    "Arcane Missiles Activation",
		Duration: core.NeverExpires,
		OnReset: func(aura *core.Aura, sim *core.Simulation) {
			aura.Activate(sim)
		},
		OnCastComplete: func(aura *core.Aura, sim *core.Simulation, spell *core.Spell) {
			if !spell.ProcMask.Matches(core.ProcMaskSpellDamage) {
				return
			}
			if sim.Proc(procChance, "Arcane Missiles") {
				mage.ArcaneMissilesProcAura.Activate(sim)
			}
		},
	})
}

func (mage *Mage) hasChillEffect(spell *core.Spell) bool {
	return spell == mage.Frostbolt || spell == mage.FrostfireBolt || (spell == mage.Blizzard && mage.Talents.IceShards > 0)
}

const (
	MageSpellFlagNone      int64 = 0
	MageSpellArcaneBarrage int64 = 1 << iota
	MageSpellArcaneBlast
	MageSpellArcaneExplosion
	MageSpellArcanePower
	MageSpellArcaneMissilesCast
	MageSpellArcaneMissilesTick
	MageSpellBlastWave
	MageSpellBlizzard
	MageSpellCombustion
	MageSpellConeOfCold
	MageSpellDeepFreeze
	MageSpellDragonsBreath
	MageSpellEvocation
	MageSpellFireBlast
	MageSpellFireball
	MageSpellFlamestrike
	MageSpellFlameOrb
	MageSpellFocusMagic
	MageSpellFreeze
	MageSpellFrostbolt
	MageSpellFrostfireBolt
	MageSpellFrostfireOrb
	MageSpellIceLance
	MageSpellIcyVeins
	MageSpellIgnite
	MageSpellLivingBombExplosion
	MageSpellLivingBombDot
	MageSpellManaGems
	MageSpellMirrorImage
	MageSpellPresenceOfMind
	MageSpellPyroblast
	MageSpellPyroblastDot
	MageSpellScorch
	MageSpellMoltenArmor
	MageSpellMageArmor

	MageSpellLast
	MageSpellsAll        = MageSpellLast<<1 - 1
	MageSpellLivingBomb  = MageSpellLivingBombDot | MageSpellLivingBombExplosion
	MageSpellFireMastery = MageSpellLivingBombDot | MageSpellPyroblastDot | MageSpellCombustion // Ignite done manually in spell due to unique mechanic
	MageSpellFire        = MageSpellBlastWave | MageSpellCombustion | MageSpellDragonsBreath | MageSpellFireball |
		MageSpellFireBlast | MageSpellFlameOrb | MageSpellFlamestrike | MageSpellFrostfireBolt | MageSpellIgnite |
		MageSpellLivingBomb | MageSpellPyroblast | MageSpellScorch
	MageSpellChill        = MageSpellFrostbolt | MageSpellFrostfireBolt
	MageSpellBrainFreeze  = MageSpellFireball | MageSpellFrostfireBolt
	MageSpellsAllDamaging = MageSpellArcaneBarrage | MageSpellArcaneBlast | MageSpellArcaneExplosion | MageSpellArcaneMissilesTick | MageSpellBlastWave | MageSpellBlizzard | MageSpellDeepFreeze |
		MageSpellDragonsBreath | MageSpellFireBlast | MageSpellFireball | MageSpellFlamestrike | MageSpellFlameOrb | MageSpellFrostbolt | MageSpellFrostfireBolt |
		MageSpellFrostfireOrb | MageSpellIceLance | MageSpellLivingBombExplosion | MageSpellLivingBombDot | MageSpellPyroblast | MageSpellPyroblastDot | MageSpellScorch
	MageSpellInstantCast = MageSpellArcaneBarrage | MageSpellArcaneMissilesCast | MageSpellArcaneMissilesTick | MageSpellFireBlast | MageSpellArcaneExplosion |
		MageSpellBlastWave | MageSpellCombustion | MageSpellConeOfCold | MageSpellDeepFreeze | MageSpellDragonsBreath | MageSpellIceLance |
		MageSpellManaGems | MageSpellMirrorImage | MageSpellPresenceOfMind | MageSpellMoltenArmor | MageSpellMageArmor
	MageSpellExtraResult = MageSpellLivingBombExplosion | MageSpellArcaneMissilesTick | MageSpellBlizzard
)
