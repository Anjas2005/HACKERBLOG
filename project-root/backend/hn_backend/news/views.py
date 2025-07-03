from rest_framework.decorators import api_view
from rest_framework.response import Response
from rest_framework import status
from .serializers import top_news_serializer
from .models import top_news

@api_view(['POST'])
def add_top_news(request):
  print("[DEBUG] Incoming data:", request.data)  # 👈 Add this line
  if isinstance(request.data, list):
        print("[INFO] Handling a bulk insert")
        serializer = top_news_serializer(data=request.data, many=True)
  else:
        print("[INFO] Handling a single insert")
        serializer = top_news_serializer(data=request.data)

  if serializer.is_valid():
      serializer.save()
      return Response(serializer.data, status=status.HTTP_201_CREATED)
  else:
      print("[ERROR] Validation failed:", serializer.errors)
      return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)


# @api_view(['POST'])
# def add_top_news(request):
#     print("[DEBUG] Incoming data type:", type(request.data))
#     print("[DEBUG] Incoming data:", request.data[:1])  # Show a sample for confirmation

#     if isinstance(request.data, list):
#         serializer = top_news_serializer(data=request.data, many=True)
#     else:
#         serializer = top_news_serializer(data=request.data)

#     if serializer.is_valid():
#         serializer.save()
#         return Response(serializer.data, status=status.HTTP_201_CREATED)
#     else:
#         print("[ERROR] Validation failed:", serializer.errors)
#         return Response(serializer.errors, status=status.HTTP_400_BAD_REQUEST)


