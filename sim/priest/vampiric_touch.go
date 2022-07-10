package priest

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (priest *Priest) registerVampiricTouchSpell() {
	actionID := core.ActionID{SpellID: 48160}
	baseCost := priest.BaseMana() * 0.16

	applier := priest.OutcomeFuncTick()
	if priest.Talents.Shadowform {
		applier = priest.OutcomeFuncMagicCrit(priest.SpellCritMultiplier(1, 1))
	}

	priest.VampiricTouch = priest.RegisterSpell(core.SpellConfig{
		ActionID:    actionID,
		SpellSchool: core.SpellSchoolShadow,

		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * 1500,
			},
		},

		ApplyEffects: core.ApplyEffectFuncDirectDamage(core.SpellEffect{
			ProcMask:            core.ProcMaskSpellDamage,
			BonusSpellHitRating: float64(priest.Talents.ShadowFocus) * 1 * core.SpellHitRatingPerHitChance,
			ThreatMultiplier:    1 - 0.08*float64(priest.Talents.ShadowAffinity),
			OutcomeApplier:      priest.OutcomeFuncMagicHit(),
			OnSpellHitDealt: func(sim *core.Simulation, spell *core.Spell, spellEffect *core.SpellEffect) {
				if spellEffect.Landed() {
					priest.VampiricTouchDot.Apply(sim)
				}
			},
		}),
	})

	target := priest.CurrentTarget

	priest.VampiricTouchDot = core.NewDot(core.Dot{
		Spell: priest.VampiricTouch,
		Aura: target.RegisterAura(core.Aura{
			Label:    "VampiricTouch-" + strconv.Itoa(int(priest.Index)),
			ActionID: actionID,
		}),

		NumberOfTicks:       5,
		TickLength:          time.Second * 3,
		AffectedByCastSpeed: priest.Talents.Shadowform,

		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			DamageMultiplier:     1 + float64(priest.Talents.Darkness)*0.02,
			BonusSpellCritRating: float64(priest.Talents.MindMelt) * 3 * core.CritRatingPerCritChance,
			ThreatMultiplier:     1 - 0.08*float64(priest.Talents.ShadowAffinity),
			IsPeriodic:           true,
			ProcMask:             core.ProcMaskPeriodicDamage,
			BaseDamage:           core.BaseDamageConfigMagicNoRoll(850/5, 0.4),
			OutcomeApplier:       applier,
		}),
	})
}
