namespace CoreService.Service.Interfaces
{
    public interface IIntentTagService
    {
        string[] GetTagsByIntent(string intent);
    }
}
