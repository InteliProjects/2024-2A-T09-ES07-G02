using CoreService.Service.Interfaces;

namespace CoreService.Service.Services
{
    public class IntentTagService : IIntentTagService
    {
        private static readonly Dictionary<string, string[]> intentTags = new Dictionary<string, string[]>
        {
            { "intent_investment", new[] { "Investimento" } },
            { "intent_stock_exchange", new[] { "Bolsa" } },
            { "intent_energy", new[] { "Energia" } },
            { "intent_macroeconomics", new[] { "Macroeconomia" } },
            { "intent_external_relations", new[] { "Relações Externas" } },
            { "intent_sustainability", new[] { "Sustentabilidade" } }
        };

        public string[] GetTagsByIntent(string intent)
        {
            return intentTags.TryGetValue(intent, out var tags) ? tags : Array.Empty<string>();
        }
    }
}