IF NOT EXISTS (SELECT name FROM sys.databases WHERE name = N'find-my-friend')
BEGIN
    CREATE DATABASE [find-my-friend];
END
GO