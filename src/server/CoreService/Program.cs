using CoreService.Application.Interfaces.Messaging;
using CoreService.Infra.Data.Context;
using CoreService.Infra.Data.Repository.Interfaces;
using CoreService.Infra.Data.Repository.Repositories;
using CoreService.Infra.Messaging;
using CoreService.Service.Interfaces;
using CoreService.Service.Services;
using Microsoft.EntityFrameworkCore;
using Newtonsoft.Json.Serialization;

var builder = WebApplication.CreateBuilder(args);

// enable Newtonsoft.Json
builder.Services.AddControllers()
    .AddNewtonsoftJson(options =>
    {
        options.SerializerSettings.ContractResolver = new CamelCasePropertyNamesContractResolver();
    });

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
builder.Services.AddAutoMapper(typeof(Program).Assembly);

// Connection String
builder.Services.AddDbContext<DatabaseContext>((serviceProvider, options) =>
{
    var configuration = serviceProvider.GetRequiredService<IConfiguration>();
    object value = options.UseNpgsql(configuration.GetConnectionString("DefaultConnection"));
});

// Kafka
builder.Services.AddScoped<IKafkaConsumer, KafkaConsumer>();
builder.Services.AddScoped<IKafkaProducer, KafkaProducer>();

// Tag
builder.Services.AddScoped<ITagRepository, TagRepository>();
builder.Services.AddScoped<ITagService, TagService>();


// Lrr
builder.Services.AddScoped<ILrrRepository, LrrRepository>();
builder.Services.AddScoped<ILrrService, LrrService>();

// LrrTag
builder.Services.AddScoped<ILrrTagRepository, LrrTagRepository>();
builder.Services.AddScoped<ILrrTagService, LrrTagService>();

// IntentTag
builder.Services.AddScoped<IIntentTagService, IntentTagService>();

// RegistredAddress
builder.Services.AddScoped<IRegistredAddressRepository, RegistredAddressRepository>();
builder.Services.AddScoped<IRegistredAddressService, RegistredAddressService>();

// enable CORS to communicate between services 
builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowAllOrigins", builder =>
    {
        builder.AllowAnyOrigin()
               .AllowAnyMethod()
               .AllowAnyHeader();
    });

    options.AddPolicy("AllowFrontend",
        policy =>
        {
            policy.WithOrigins("http://localhost:5173")
                  .AllowAnyHeader()
                  .AllowAnyMethod();
        });
});

var app = builder.Build();

// pipe for HTTP requests
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseRouting();

// CORS implementation
app.UseCors("AllowAllOrigins");

app.MapControllers();

app.Run();
