FROM 1.19.7-alpine3.17

LABEL maintainer="hans.kusos@hotmail.com"

###############################################################################
#                                INSTALLATION
###############################################################################

# 設置固定的項目路徑
ENV WORKDIR /usr/src/app

# 添加應用可執行文件，並設置執行權限
COPY .   $WORKDIR/main
RUN chmod +x $WORKDIR/main

# 添加I18N多語言文件、靜態文件、配置文件、模板文件
ADD i18n     $WORKDIR/i18n
ADD public   $WORKDIR/public
ADD config   $WORKDIR/config
ADD template $WORKDIR/template

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./main
