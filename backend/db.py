import mysql.connector
from mysql.connector import Error

def init_db():
    try:
        # Connect to MySQL server
        connection = mysql.connector.connect(
            host='your_host',
            user='your_user',
            password='your_password'
        )

        if connection.is_connected():
            cursor = connection.cursor()

            # Create database
            cursor.execute("CREATE DATABASE IF NOT EXISTS your_database")

            # Switch to the newly created database
            cursor.execute("USE your_database")

            # Create table
            cursor.execute("""
                CREATE TABLE IF NOT EXISTS Feedback (
                    id INT AUTO_INCREMENT PRIMARY KEY,
                    answer VARCHAR(255) NOT NULL,
                    suggestion TEXT
                )
            """)

            print("Database initialized successfully.")

    except Error as e:
        print("Error:", e)

    finally:
        # Close connection
        if connection.is_connected():
            cursor.close()
            connection.close()

if __name__ == "__main__":
    init_db()
