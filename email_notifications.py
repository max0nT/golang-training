import pika
import smtplib
import logging
from email.mime.text import MIMEText

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

def send_email_smtp(message_text: str):
    try:
        msg = MIMEText(message_text)
        msg['Subject'] = 'Уведомление из RabbitMQ'
        msg['From'] = 'sender@example.com'
        msg['To'] = 'recipient@example.com'

        with smtplib.SMTP('mailpit', 1025) as server:
            server.send_message(msg)
            logger.info("Email отправлен через Mailpit")

    except Exception as e:
        logger.error(f"Ошибка отправки email: {e}")

def process_message(ch, method, properties, body):
    try:
        message_text = body.decode('utf-8')
        logger.info(f"Получено сообщение: {message_text}")

        send_email_smtp(message_text)

        ch.basic_ack(delivery_tag=method.delivery_tag)

    except Exception as e:
        logger.error(f"Ошибка обработки: {e}")
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

def main():
    try:
        connection = pika.BlockingConnection(
            pika.ConnectionParameters('0.0.0.0', 5672)
        )
        channel = connection.channel()

        channel.queue_declare(queue='email_notifications', durable=True)

        channel.basic_consume(
            queue='email_notifications',
            on_message_callback=process_message,
            auto_ack=False
        )

        logger.info("Ожидание сообщений...")
        channel.start_consuming()

    except Exception as e:
        logger.error(f"Ошибка: {e}")

if __name__ == '__main__':
    main()
