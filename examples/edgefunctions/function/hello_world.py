import os

def handler(req, context):
    name = os.getenv('NAME') or 'Unknown'

    return {
        'body': f'Hello, {name}!',
        'statusCode': 200
    }
